package utility

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/andybalholm/brotli"
)

func HeaderToString(header map[string][]string) string {
	txt := ""
	for key, values := range header {
		for _, value := range values {
			txt = txt + key + ": " + value + "\n"
		}
	}
	return txt
}

func DecodeContent(encoding string, buf *bytes.Buffer) (*[]byte, error) {
	var reader io.Reader
	var err error

	switch strings.ToLower(encoding) {
	case "br":
		reader = brotli.NewReader(buf)

	case "gzip":
		reader, err = gzip.NewReader(buf)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %v", err)
		}

	case "deflate":
		reader = flate.NewReader(buf)

	default:
		content, err := io.ReadAll(buf)
		if err != nil {
			return nil, err
		}
		return &content, nil
	}

	decodedBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read decoded content: %v", err)
	}

	if r, ok := reader.(io.Closer); ok {
		r.Close()
	}

	return &decodedBytes, nil
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	if unicode.IsLower(r[0]) {
		r[0] = unicode.ToUpper(r[0])
	}

	return string(r)
}
