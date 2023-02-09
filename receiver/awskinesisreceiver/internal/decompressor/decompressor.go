package decompressor

import (
	"bytes"
	"compress/gzip"
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
	format       string
}

func NewDecompressor(format string) DeCompressor {
	return &decompressor{format: format}
}

func (d *decompressor) reader(buffer *bytes.Buffer, format string) (bufferedResetReader, error) {
	switch format {
	case "gzip":
		r, err := gzip.NewReader(buffer)
		if err != nil {
			return nil, err
		}
		return r, nil
	case "flate":
		r, err := NewInflateReader(buffer)
		if err != nil {
			return nil, err
		}
		return r, nil
	case "zlib":
		r, err := NewZlibReader(buffer)
		if err != nil {
			return nil, err
		}
		return r, nil
	case "noop", "none":
		return NewReader(buffer), nil
	default:
		return NewReader(buffer), nil
	}
}
func (d *decompressor) Do(in []byte) ([]byte, error) {
	inBuf := new(bytes.Buffer)
	inBuf.Write(in)
	if d.decompressor == nil {
		decompressor, err := d.reader(inBuf, d.format)
		if err != nil {
			return nil, err
		}
		d.decompressor = decompressor
	} else {
		if err := d.decompressor.Reset(inBuf); err != nil {
			return nil, err
		}
	}
	outBuf := new(bytes.Buffer)
	_, err := io.Copy(outBuf, d.decompressor)
	if err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}
