package factory

import (
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
	"github.com/erashaperavm/BartmossRoar/core/crypto/blake3"
	"github.com/erashaperavm/BartmossRoar/core/crypto/sha"
)

func NewHasher(algo crypto.HashAlgo) (crypto.Hasher, error) {
	switch algo {
	case crypto.HashBlake3:
		return blake3.New(), nil
	case crypto.HashSHA256:
		return sha.New(), nil
	default:
		return nil, fmt.Errorf("unknown hash algo: %d", algo)
	}
}
