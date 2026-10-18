package factory

import (
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
	"github.com/erashaperavm/BartmossRoar/core/crypto/circl_sign/ed25519"
	"github.com/erashaperavm/BartmossRoar/core/crypto/circl_sign/mldsa"
)

// NewSigner 负责把具体的实现返回给上层。
func NewSigner(algo crypto.SignAlgo) (crypto.Signer, crypto.PublicKey, error) {
	switch algo {
	case crypto.SignMLDSA44:
		return mldsa.New44()
	case crypto.SignMLDSA65:
		return mldsa.New65()
	case crypto.SignMLDSA87:
		return mldsa.New87()
	case crypto.SignEd25519:
		return ed25519.New()
	default:
		return nil, nil, fmt.Errorf("unknown signature algo: %d", algo)
	}
}
