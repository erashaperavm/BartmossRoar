package transport

type Proto byte
type Addr string

const (
	TCP Proto = iota
	Tor
)

type NodeAddr struct {
	Protocol Proto
	Address  Addr
}

type Handler = func(from NodeAddr, data []byte)

// Transport 只负责选择协议、收发字节
type Transport interface {
	Send(addr NodeAddr, data []byte) error
	OnRecv(handler Handler)
	LocalAddr() NodeAddr
	Close() error
}
