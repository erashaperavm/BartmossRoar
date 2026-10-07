package transport

import (
	"errors"
	"fmt"
	"sync"
)

type Node struct {
	Trans    Transport
	Mu       sync.Mutex
	Protocol Proto
}

func NewNode(protocol Proto, addrIfTcp string, localPort int, connectExistingIfTorHS bool) (*Node, error) {
	switch protocol {
	case TCP:
		trans, err := NewTCP(addrIfTcp)
		if err != nil {
			return nil, err
		}
		return &Node{
			Trans:    trans,
			Mu:       sync.Mutex{},
			Protocol: protocol,
		}, nil
	case Tor:
		trans, err := NewTor(localPort, connectExistingIfTorHS)
		if err != nil {
			return nil, err
		}
		return &Node{
			Trans:    trans,
			Mu:       sync.Mutex{},
			Protocol: protocol,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("unknown protocol: %v", protocol))
}
