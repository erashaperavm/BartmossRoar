package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/nao1215/tornago"
)

type TorTransport struct {
	ctx    context.Context
	cancel context.CancelFunc

	ln net.Listener

	mu      sync.RWMutex
	conns   map[Addr]net.Conn
	closers map[Addr]io.Closer // 与 conns 一一对应，用于释放 tor client 等附属资源

	handlerMu sync.RWMutex
	handler   Handler

	onionAddr string
	socksAddr string

	tor        *tornago.TorProcess
	ctrlClient *tornago.ControlClient
	hiddenSvc  tornago.HiddenService
	ownTor     bool
	closeOnce  sync.Once
}

// NewTor 创建 Tor 传输层。
// connectExisting == false:
//
//	启动一个临时 Tor 实例。
//
// connectExisting == true:
//
//	连接已经运行的 Tor：
//	SOCKS 127.0.0.1:9050
//	Control 127.0.0.1:9051
func NewTor(localPort int, connectExisting bool) (*TorTransport, error) {
	ctx, cancel := context.WithCancel(context.Background())

	t := &TorTransport{
		ctx:     ctx,
		cancel:  cancel,
		conns:   make(map[Addr]net.Conn),
		closers: make(map[Addr]io.Closer),
	}

	// 1. 获取 / 启动 Tor
	var (
		tor *tornago.TorProcess
		err error
	)

	if connectExisting {
		t.socksAddr = "127.0.0.1:9050"
	} else {
		launchCfg, err := tornago.NewTorLaunchConfig(
			tornago.WithTorSocksAddr(":0"),
			tornago.WithTorControlAddr(":0"),
			tornago.WithTorStartupTimeout(60*time.Second),
		)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("create tor launch config: %w", err)
		}

		tor, err = tornago.StartTorDaemon(launchCfg)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("start tor daemon: %w", err)
		}

		t.tor = tor
		t.ownTor = true
		t.socksAddr = tor.SocksAddr()
	}

	// 2. 建立 ControlPort
	controlAddr := "127.0.0.1:9051"
	if tor != nil {
		controlAddr = tor.ControlAddr()
	}

	auth, _, err := tornago.ControlAuthFromTor(controlAddr, 30*time.Second)
	if err != nil {
		t.cleanup()
		return nil, fmt.Errorf("tor control authentication: %w", err)
	}

	ctrlClient, err := tornago.NewControlClient(controlAddr, auth, 30*time.Second)
	if err != nil {
		t.cleanup()
		return nil, fmt.Errorf("create control client: %w", err)
	}
	if err := ctrlClient.Authenticate(); err != nil {
		ctrlClient.Close()
		t.cleanup()
		return nil, fmt.Errorf("authenticate control client: %w", err)
	}
	t.ctrlClient = ctrlClient

	// 3. 先监听本地 TCP
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		t.cleanup()
		return nil, fmt.Errorf("local listen: %w", err)
	}
	t.ln = ln

	// 4. 创建 Hidden Service
	hsCfg, err := tornago.NewHiddenServiceConfig(
		tornago.WithHiddenServicePort(80, localPort),
	)
	if err != nil {
		t.cleanup()
		return nil, fmt.Errorf("create hidden service config: %w", err)
	}

	hs, err := ctrlClient.CreateHiddenService(ctx, hsCfg)
	if err != nil {
		t.cleanup()
		return nil, fmt.Errorf("create hidden service: %w", err)
	}
	t.hiddenSvc = hs
	t.onionAddr = hs.OnionAddress()

	go t.acceptLoop()
	return t, nil
}

func (t *TorTransport) acceptLoop() {
	backoff := 5 * time.Millisecond
	const maxBackoff = time.Second

	for {
		conn, err := t.ln.Accept()
		if err != nil {
			select {
			case <-t.ctx.Done():
				return
			default:
			}
			if errors.Is(err, net.ErrClosed) {
				return
			}
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}
		backoff = 5 * time.Millisecond
		go t.handleConn(conn)
	}
}

// handleConn 处理入站连接。
//
// 注意：Hidden Service 入站连接的 conn.RemoteAddr() 是本地 Tor
// 代理地址，不是对方的真实 .onion 地址。上层协议必须从 payload
// 中提取 FromID 才能知道对方是谁。这里的 Address 仅作占位。
func (t *TorTransport) handleConn(conn net.Conn) {
	remote := NodeAddr{
		Protocol: Tor,
		Address:  Addr(conn.RemoteAddr().String()),
	}
	t.handleStream(remote, conn)
}

// handleStream 是统一的读循环，入站和出站共用。
func (t *TorTransport) handleStream(remote NodeAddr, conn net.Conn) {
	defer func() {
		conn.Close()
		t.removeConn(remote.Address, conn)
	}()

	// 同 tcp.go：不做分帧，分帧属于协议层。
	buf := make([]byte, 64*1024)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			t.dispatch(remote, data)
		}
		if err != nil {
			return
		}
	}
}

func (t *TorTransport) dispatch(from NodeAddr, data []byte) {
	t.handlerMu.RLock()
	h := t.handler
	t.handlerMu.RUnlock()
	if h != nil {
		h(from, data)
	}
}

// Send 通过 Tor SOCKS5 向目标 Onion Service 发送数据。
// 优先复用 conns 里的连接，没有才新建 tor client 拨号。
func (t *TorTransport) Send(addr NodeAddr, data []byte) error {
	if addr.Protocol != Tor {
		return fmt.Errorf("invalid protocol: expected Tor, got %v", addr.Protocol)
	}

	onion := normalizeOnionAddr(string(addr.Address))
	if onion == "" {
		return errors.New("empty onion address")
	}
	if !strings.HasSuffix(onion, ".onion") {
		return fmt.Errorf("not an onion address: %s", onion)
	}

	t.mu.RLock()
	conn, ok := t.conns[addr.Address]
	t.mu.RUnlock()

	if !ok {
		c, closer, err := t.dialOnion(onion)
		if err != nil {
			return err
		}

		t.mu.Lock()
		if existing, dup := t.conns[addr.Address]; dup {
			// 竞态：别人抢先了，丢弃自己这份
			t.mu.Unlock()
			c.Close()
			if closer != nil {
				closer.Close()
			}
			conn = existing
		} else {
			t.conns[addr.Address] = c
			if closer != nil {
				t.closers[addr.Address] = closer
			}
			t.mu.Unlock()
			conn = c
			// 出站连接启动读循环，接收对方回包
			go t.handleStream(addr, c)
		}
	}

	if _, err := conn.Write(data); err != nil {
		// 主动关闭，让读循环退出并触发清理
		conn.Close()
		return fmt.Errorf("tor write %s: %w", onion, err)
	}
	return nil
}

// dialOnion 建立一条到目标 onion 的连接。
// 返回的 closer 是 tor client，必须在 conn 生命周期结束时一起关闭。
func (t *TorTransport) dialOnion(onion string) (net.Conn, io.Closer, error) {
	clientCfg, err := tornago.NewClientConfig(
		tornago.WithClientSocksAddr(t.socksAddr),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create tor client config: %w", err)
	}
	client, err := tornago.NewClient(clientCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("create tor client: %w", err)
	}

	target := net.JoinHostPort(onion, "80")
	conn, err := client.Dial("tcp", target)
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("tor dial %s: %w", target, err)
	}

	// 不 defer client.Close()：conn 的生命周期独立于本次 Send。
	// client 与 conn 绑定，一起存到 closers，等 conn 关闭时再释放。
	return conn, client, nil
}

// normalizeOnionAddr 接受以下形式：
//
//	abcdef.onion
//	abcdef.onion:80
//	http://abcdef.onion
//	http://abcdef.onion:80
//
// 最终统一返回 hostname。
func normalizeOnionAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	addr = strings.TrimPrefix(addr, "tcp://")
	addr = strings.TrimPrefix(addr, "http://")
	addr = strings.TrimPrefix(addr, "https://")
	addr = strings.TrimSuffix(addr, "/")

	if host, _, err := net.SplitHostPort(addr); err == nil {
		addr = host
	}
	return addr
}

func (t *TorTransport) removeConn(addr Addr, conn net.Conn) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if existing, ok := t.conns[addr]; ok && existing == conn {
		delete(t.conns, addr)
		// 释放 conn 绑定的附属资源（tor client）
		if c, hasCloser := t.closers[addr]; hasCloser {
			c.Close()
			delete(t.closers, addr)
		}
	}
}

func (t *TorTransport) OnRecv(handler Handler) {
	t.handlerMu.Lock()
	t.handler = handler
	t.handlerMu.Unlock()
}

func (t *TorTransport) LocalAddr() NodeAddr {
	return NodeAddr{
		Protocol: Tor,
		Address:  Addr(t.onionAddr),
	}
}

func (t *TorTransport) Close() error {
	var closeErr error

	t.closeOnce.Do(func() {
		t.cancel()

		if t.ln != nil {
			if err := t.ln.Close(); err != nil {
				closeErr = err
			}
		}

		// 删除 Onion Service，给个短超时
		if t.hiddenSvc != nil {
			rmCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := t.hiddenSvc.Remove(rmCtx); err != nil && closeErr == nil {
				closeErr = err
			}
			cancel()
		}

		// 关闭所有 conn 及其绑定的 closer
		t.mu.Lock()
		for addr, conn := range t.conns {
			_ = conn.Close()
			delete(t.conns, addr)
		}
		for addr, c := range t.closers {
			_ = c.Close()
			delete(t.closers, addr)
		}
		t.mu.Unlock()

		if t.ctrlClient != nil {
			if err := t.ctrlClient.Close(); err != nil && closeErr == nil {
				closeErr = err
			}
		}

		if t.tor != nil && t.ownTor {
			if err := t.tor.Stop(); err != nil && closeErr == nil {
				closeErr = err
			}
		}
	})

	return closeErr
}

// cleanup 用于 NewTor 失败时回滚已创建资源。
func (t *TorTransport) cleanup() {
	if t.hiddenSvc != nil {
		rmCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = t.hiddenSvc.Remove(rmCtx)
		cancel()
	}
	if t.ln != nil {
		_ = t.ln.Close()
	}
	if t.ctrlClient != nil {
		_ = t.ctrlClient.Close()
	}
	if t.tor != nil && t.ownTor {
		_ = t.tor.Stop()
	}
	t.cancel()
}
