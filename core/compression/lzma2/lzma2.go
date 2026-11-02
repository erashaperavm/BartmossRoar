package lzma2

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/erashaperavm/BartmossRoar/core/compression"
	"github.com/ulikunitz/xz/lzma"
)

type LZMA2Compressor struct {
}

func (c *LZMA2Compressor) CompressStream(in io.Reader, out io.Writer) error {
	if in == nil || out == nil {
		return errors.New("input or output stream is nil")
	}

	w, err := lzma.NewWriter2(out)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, in)
	if err != nil {
		return err
	}

	return nil
}

// CompressManyStream supports concurrent compression for each stream, it will compress each origin concurrently.
func (c *LZMA2Compressor) CompressManyStream(in []io.Reader, out []io.Writer) []error {
	if len(in) != len(out) {
		return []error{errors.New("input and output stream slices must have the same length")}
	}

	errs := make([]error, len(in))
	var wg sync.WaitGroup

	for i := range in {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if err := c.CompressStream(in[idx], out[idx]); err != nil {
				errs[idx] = err
			}
		}(i)
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return errs
		}
	}
	return nil
}

func (c *LZMA2Compressor) ExtractStream(in io.Reader, out io.Writer) error {
	if in == nil || out == nil {
		return errors.New("input or output stream is nil")
	}

	r, err := lzma.NewReader2(in)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	return nil
}

// ExtractManyStream supports concurrent extracting, it will extract each origin concurrently.
func (c *LZMA2Compressor) ExtractManyStream(in []io.Reader, out []io.Writer) error {
	if len(in) != len(out) {
		return errors.New("input and output stream slices must have the same length")
	}

	var wg sync.WaitGroup
	errs := make([]error, len(in))

	for i := range in {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if err := c.ExtractStream(in[idx], out[idx]); err != nil {
				errs[idx] = err
			}
		}(i)
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("one or more errors occurred during extraction: %v", errs)
		}
	}
	return nil
}

func (c *LZMA2Compressor) Algorithm() compression.CompressAlgo {
	return compression.CompressLZMA2
}

func (c *LZMA2Compressor) Compress(origin []byte) ([]byte, error) {
	return nil, errors.New("lzma2 does not support Compress method, please use CompressStream or CompressManyStream")
}

func (c *LZMA2Compressor) CompressMany(origins [][]byte) ([][]byte, []error) {
	return nil, []error{errors.New("lzma2 does not support CompressMany method, please use CompressStream or CompressManyStream")}
}

func (c *LZMA2Compressor) Extract(compressed []byte) ([]byte, error) {
	return nil, errors.New("lzma2 does not support Extract method, please use ExtractStream or ExtractManyStream")
}

func (c *LZMA2Compressor) ExtractMany(compresseds [][]byte) ([][]byte, []error) {
	return nil, []error{errors.New("lzma2 does not support ExtractMany method, please use ExtractStream or ExtractManyStream")}
}
