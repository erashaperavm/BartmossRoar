package factory

import (
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
	"github.com/erashaperavm/BartmossRoar/core/crypto/hybrid"
)

func NewKEM(algo crypto.KemAlgo) (crypto.KEMScheme, error) {
	switch algo {
	case crypto.X25519MLKEM768:
		return hybrid.NewScheme(), nil
	default:
		return nil, fmt.Errorf("unknown symmetric algo: %d", algo)
	}
}
