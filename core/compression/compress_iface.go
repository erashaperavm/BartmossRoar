package compression

import "io"

type CompressAlgo byte

const (
	CompressLZMA2 CompressAlgo = iota
	CompressZSTD
)

type Compressor interface {
	CompressStream(in io.Reader, out io.Writer) error
	CompressManyStream(in []io.Reader, out []io.Writer) []error
	ExtractStream(in io.Reader, out io.Writer) error
	ExtractManyStream(in []io.Reader, out []io.Writer) error
	Compress(origin []byte) ([]byte, error)
	CompressMany(origins [][]byte) ([][]byte, []error)
	Extract(compressed []byte) ([]byte, error)
	ExtractMany(compresseds [][]byte) ([][]byte, []error)
	Algorithm() CompressAlgo
}
