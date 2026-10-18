package factory

import (
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
	"github.com/erashaperavm/BartmossRoar/core/crypto/chacha20poly1305"
)

// NewSymmetricCipher 负责把具体的实现返回给上层。
func NewSymmetricCipher(algo crypto.EncAlgo, key crypto.Key) (crypto.SymmetricCipher, error) {
	switch algo {
	case crypto.EncChaCha20Poly1305:
		return chacha20poly1305.New(key)
	case crypto.EncXChaCha20Poly1305:
		return chacha20poly1305.NewX(key)
	default:
		return nil, fmt.Errorf("unknown symmetric algo: %d", algo)
	}
}
