package decompressor

import "io"

type noop struct {
	data io.Reader
}

func NewNoopCompressor() DeCompressor {
	return &decompressor{decompressor: &noop{}}
}

func (n *noop) Read(p []byte) (int, error) {
	return len(p), nil
}

func (n *noop) Reset(r io.Reader) error {
	n.data = r
	return nil
}
