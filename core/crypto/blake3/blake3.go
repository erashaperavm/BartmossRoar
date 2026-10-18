package blake3

import (
	"errors"
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
)

import "github.com/zeebo/blake3"

/*
	use zeebo/blake3（github.com/zeebo/blake3）
*/

type Blake3hasher struct {
}

func (b3 *Blake3hasher) Hash(b []byte) crypto.HashSum {
	return blake3.Sum256(b)
}

func (b3 *Blake3hasher) HashMany(manyBytes [][]byte) ([]crypto.HashSum, error) {
	out := make([]crypto.HashSum, len(manyBytes))
	// 复用一个 Hasher，避免每次新建
	hh := blake3.New()
	for i, b := range manyBytes {
		hh.Reset() // 关键：复用而非新建
		j, err := hh.Write(b)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%d blake3 hash error: %v", j, err))
		}
		var sum crypto.HashSum
		// Hasher.Sum(nil) 返回 []byte，长度 32
		copy(sum[:], hh.Sum(nil))
		out[i] = sum
	}

	return out, nil
}

func (b3 *Blake3hasher) Algorithm() crypto.HashAlgo {
	return crypto.HashBlake3
}

func New() *Blake3hasher {
	return &Blake3hasher{}
}
