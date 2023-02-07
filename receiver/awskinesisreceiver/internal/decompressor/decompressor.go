package decompressor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

type bufferedResetReader interface {
	Read(p []byte) (n int, err error)
	Reset(r io.Reader) error
}

type DeCompressor interface {
	Do(in []byte) (out []byte, err error)
}

var _ DeCompressor = (*decompressor)(nil)

type decompressor struct {
	decompressor bufferedResetReader
}

func NewDecompressor(format string) (DeCompressor, error) {
	d := &decompressor{
		decompressor: &noop{},
	}
	switch format {
	case "gzip":
		r, err := gzip.NewReader(nil)
		if err != nil {
			return nil, err
		}
		d.decompressor = r
	case "noop", "none":
	default:
		return nil, fmt.Errorf("unknown compression format: %s", format)
	}
	return d, nil
}

func (d *decompressor) Do(in []byte) (out []byte, err error) {
	buf := new(bytes.Buffer)
	d.decompressor.Reset(buf)
	if _, err := d.decompressor.Read(in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
