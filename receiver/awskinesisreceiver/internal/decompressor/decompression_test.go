package decompressor

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestDeCompressorFormats(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		format string
		data   []string
	}{
		{format: "gzip", data: []string{"You know nothing Jon Snow", "The black fox"}},
		{format: "flate", data: []string{"You know nothing Jon Snow", "The black fox"}},
		{format: "zlib", data: []string{"You know nothing Jon Snow", "The black fox"}},
		{format: "noop", data: []string{"You know nothing Jon Snow", "The black fox"}},
		{format: "none", data: []string{"You know nothing Jon Snow", "The black fox"}},
	}
	//const data = "You know nothing Jon Snow"
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("format_%s", tc.format), func(t *testing.T) {
			decompressor := NewDecompressor(tc.format)
			for _, data := range tc.data {
				switch tc.format {
				case "none", "noop":
					out, err := decompressor.Do([]byte(data))
					if err != nil {
						t.Errorf("decompression is failing")
					}
					if string(out) != data {
						t.Errorf("decompressed data is not matching")
					}
				case "gzip":
					buf := bytes.Buffer{}
					w, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
					if err != nil {
						t.Errorf("not able to create writer")
					}
					if _, err := io.Copy(w, strings.NewReader(data)); err != nil {
						t.Errorf("not able to compress the file")
					}
					if err := w.Close(); err != nil {
						t.Errorf("not able to close writer")
					}
					out, err := decompressor.Do(buf.Bytes())
					if err != nil {
						t.Errorf("decompression is failing")
					}
					if string(out) != data {
						t.Errorf("decompressed data is not matching")
					}
				case "zlib":
					buf := bytes.Buffer{}
					w, err := zlib.NewWriterLevel(&buf, zlib.BestSpeed)
					if err != nil {
						t.Errorf("not able to create writer")
					}
					if _, err := io.Copy(w, strings.NewReader(data)); err != nil {
						t.Errorf("not able to compress the file")
					}
					if err := w.Close(); err != nil {
						t.Errorf("not able to close writer")
					}
					out, err := decompressor.Do(buf.Bytes())
					if err != nil {
						t.Errorf("decompression is failing")
					}
					if string(out) != data {
						t.Errorf("decompressed data is not matching")
					}
				case "flate":
					buf := bytes.Buffer{}
					w, err := flate.NewWriter(&buf, flate.BestSpeed)
					if err != nil {
						t.Errorf("not able to create writer")
					}
					if _, err := io.Copy(w, strings.NewReader(data)); err != nil {
						t.Errorf("not able to compress the file")
					}
					if err := w.Close(); err != nil {
						t.Errorf("not able to close writer")
					}
					out, err := decompressor.Do(buf.Bytes())
					if err != nil {
						t.Errorf("decompression is failing")
					}
					if string(out) != data {
						t.Errorf("decompressed data is not matching")
					}
				default:
					t.Errorf("invalid format")
				}
			}

		})
	}
}
