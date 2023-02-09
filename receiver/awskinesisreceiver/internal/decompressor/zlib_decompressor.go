package decompressor

import (
	"compress/zlib"
	"fmt"
	"io"
)

type ZlibReader struct {
	zlib io.Reader
}

func NewZlibReader(r io.Reader) (*ZlibReader, error) {
	zlib, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &ZlibReader{zlib: zlib}, nil
}

func (z *ZlibReader) Read(p []byte) (int, error) {
	return z.zlib.Read(p)
}

func (z *ZlibReader) Reset(r io.Reader) error {
	reSetter, ok := z.zlib.(zlib.Resetter)
	if ok {
		return reSetter.Reset(r, nil)
	}
	return fmt.Errorf("not able to decompress the data")
}
