package factory

import (
	"github.com/erashaperavm/BartmossRoar/core/compression"
	"github.com/erashaperavm/BartmossRoar/core/compression/lzma2"
	"github.com/erashaperavm/BartmossRoar/core/compression/zstd"
)

func NewCompressor(algo compression.CompressAlgo) compression.Compressor {
	switch algo {
	case compression.CompressLZMA2:
		return &lzma2.LZMA2Compressor{}
	case compression.CompressZSTD:
		enc, _ := zstd.NewZSTDEncoder()
		dec, _ := zstd.NewZSTDDecoder()
		return &zstd.ZSTDCompressor{
			Encoder: enc,
			Decoder: dec,
		}
	default:
		return nil
	}
}
