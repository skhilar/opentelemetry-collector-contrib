package decompressor

import "io"

type noop struct {
	data io.Reader
}

func NewReader(r io.Reader) *noop {
	return &noop{data: r}
}
func (n *noop) Read(p []byte) (int, error) {
	return n.data.Read(p)
}

func (n *noop) Reset(r io.Reader) error {
	n.data = r
	return nil
}
