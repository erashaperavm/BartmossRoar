package sha

import (
	"errors"
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
)
import "crypto/sha256"

type Sha256Hasher struct {
}

func (s2 *Sha256Hasher) Hash(b []byte) crypto.HashSum {
	return sha256.Sum256(b)
}

func (s2 *Sha256Hasher) HashMany(manyBytes [][]byte) ([]crypto.HashSum, error) {
	out := make([]crypto.HashSum, len(manyBytes))
	hh := sha256.New()
	for i, b := range manyBytes {
		hh.Reset()
		j, err := hh.Write(b)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%d sha256 hash error: %v", j, err))
		}
		var sum crypto.HashSum
		// Hasher.Sum(nil) 返回 []byte，长度 32
		copy(sum[:], hh.Sum(nil))
		out[i] = sum
	}

	return out, nil
}

func (s2 *Sha256Hasher) Algorithm() crypto.HashAlgo {
	return crypto.HashSHA256
}

func New() *Sha256Hasher {
	return &Sha256Hasher{}
}
