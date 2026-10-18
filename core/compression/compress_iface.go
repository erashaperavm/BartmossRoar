package compression

type CompressAlgo byte

const (
	CompressLZMA2 CompressAlgo = iota
	CompressZSTD
)

type Compressor interface {
	Compress(origin []byte) ([]byte, error)
	CompressMany(origins [][]byte) ([][]byte, error)
	Extract(compressed []byte) ([]byte, error)
	ExtractMany(compresseds [][]byte) ([][]byte, error)
	Algorithm() CompressAlgo
}

func NewCompressor(algo CompressAlgo) (Compressor, error) {
	switch algo {

	}

	return nil, nil
}
