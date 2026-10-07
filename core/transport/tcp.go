package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type TCPTransport struct {
	ctx    context.Context
	cancel context.CancelFunc
	ln     net.Listener

	mu    sync.RWMutex
	conns map[Addr]net.Conn

	handlerMu sync.RWMutex
	handler   Handler
}

func NewTCP(addr string) (*TCPTransport, error) {
	ctx, cancel := context.WithCancel(context.Background())
	t := &TCPTransport{
		ctx:    ctx,
		cancel: cancel,
		conns:  make(map[Addr]net.Conn),
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("tcp listen: %w", err)
	}
	t.ln = ln

	go t.acceptLoop()
	return t, nil
}

func (t *TCPTransport) acceptLoop() {
	backoff := 5 * time.Millisecond
	const maxBackoff = time.Second

	for {
		conn, err := t.ln.Accept()
		if err != nil {
			// 优先判断是否主动关闭
			select {
			case <-t.ctx.Done():
				return
			default:
			}
			if errors.Is(err, net.ErrClosed) {
				return
			}
			// 避免错误下 CPU 打满：指数退避
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

// handleConn 处理入站连接：从 conn 推断 remote 地址。
func (t *TCPTransport) handleConn(conn net.Conn) {
	remote := NodeAddr{
		Protocol: TCP,
		Address:  Addr(conn.RemoteAddr().String()),
	}
	t.handleStream(remote, conn)
}

// handleStream 是入站和出站共用的读循环。
// remote 由调用方显式传入，出站时用 Send 里的 addr，
// 保证 map key 与 Send/removeConn 一致。
func (t *TCPTransport) handleStream(remote NodeAddr, conn net.Conn) {
	defer func() {
		conn.Close()
		t.removeConn(remote.Address, conn)
	}()

	// 说明：这里没有做消息分帧（length prefix / 固定块切分）。
	// conn.Read 是流式的，一次 Read 可能拿到半条或多条 block。
	// 分帧是协议层的职责，transport 只负责字节流搬运。
	// 上层需要自己维护缓冲区，按 block 长度切分后再解析。
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

func (t *TCPTransport) dispatch(from NodeAddr, data []byte) {
	t.handlerMu.RLock()
	h := t.handler
	t.handlerMu.RUnlock()
	if h != nil {
		h(from, data)
	}
}

func (t *TCPTransport) Send(addr NodeAddr, data []byte) error {
	t.mu.RLock()
	conn, ok := t.conns[addr.Address]
	t.mu.RUnlock()

	if !ok {
		c, err := net.Dial("tcp", string(addr.Address))
		if err != nil {
			return fmt.Errorf("tcp dial %s: %w", addr.Address, err)
		}

		t.mu.Lock()
		// double-check：可能其他 goroutine 抢先建立了连接
		if existing, dup := t.conns[addr.Address]; dup {
			t.mu.Unlock()
			c.Close()
			conn = existing
		} else {
			t.conns[addr.Address] = c
			t.mu.Unlock()
			conn = c
			// 关键修复：出站连接必须启动读循环，
			// 否则对方回包永远进不来 handler。
			go t.handleStream(addr, c)
		}
	}

	if _, err := conn.Write(data); err != nil {
		// 主动关闭，让读循环读到 EOF 后由 removeConn 清理。
		conn.Close()
		return fmt.Errorf("tcp write %s: %w", addr.Address, err)
	}
	return nil
}

func (t *TCPTransport) removeConn(addr Addr, conn net.Conn) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if existing, ok := t.conns[addr]; ok && existing == conn {
		delete(t.conns, addr)
	}
}

func (t *TCPTransport) OnRecv(handler Handler) {
	t.handlerMu.Lock()
	t.handler = handler
	t.handlerMu.Unlock()
}

func (t *TCPTransport) LocalAddr() NodeAddr {
	return NodeAddr{Protocol: TCP, Address: Addr(t.ln.Addr().String())}
}

func (t *TCPTransport) Close() error {
	t.cancel()
	err := t.ln.Close()

	t.mu.Lock()
	for _, conn := range t.conns {
		conn.Close()
	}
	t.conns = make(map[Addr]net.Conn)
	t.mu.Unlock()

	if err != nil && !errors.Is(err, net.ErrClosed) {
		return err
	}
	return nil
}
