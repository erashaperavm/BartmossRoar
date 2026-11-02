package factory

import (
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
	"github.com/erashaperavm/BartmossRoar/core/crypto/circl_sign/bls/BLS12_381"
)

func NewBLSSigner(algo crypto.BlsAlgo) (crypto.BlsSigner, error) {
	switch algo {
	case crypto.Bls12381:
		signer, _, err := BLS12_381.Scheme{}.GenerateKeyPair()
		if err != nil {
			return nil, err
		}
		return signer, nil
	default:
		return nil, fmt.Errorf("unknown bls algo: %d", algo)
	}
}
