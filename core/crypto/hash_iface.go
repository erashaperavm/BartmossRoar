package crypto

type HashAlgo byte

const (
	HashBlake3 HashAlgo = iota
	HashSHA256
)

// HashSum 只支持 32 字节输出算法
type HashSum [32]byte

type Hasher interface {
	Hash(b []byte) HashSum
	HashMany(manyBytes [][]byte) ([]HashSum, error)
	Algorithm() HashAlgo
}
