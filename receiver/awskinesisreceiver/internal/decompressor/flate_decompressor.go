package decompressor

import (
	"compress/flate"
	"fmt"
	"io"
)

type InflateReader struct {
	inflate io.ReadCloser
}

func NewInflateReader(r io.Reader) (*InflateReader, error) {
	inflateReader := &InflateReader{inflate: flate.NewReader(r)}
	return inflateReader, nil
}

func (f *InflateReader) Read(p []byte) (int, error) {
	return f.inflate.Read(p)
}

func (f *InflateReader) Reset(r io.Reader) error {
	reSetter, ok := f.inflate.(flate.Resetter)
	if ok {
		return reSetter.Reset(r, nil)
	}
	return fmt.Errorf("not able to uncompress the data")
}
