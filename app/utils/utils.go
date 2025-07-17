package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"slices"
)

func getSupportedEncodings() []string {
	return []string{"gzip"}
}

func IsEncodingSupported(title string) bool {
	return slices.Contains(getSupportedEncodings(), title)
}

func CompressString(s string, encoding string) ([]byte, error) {
	if encoding == "gzip" {
		var buffer bytes.Buffer
		writer := gzip.NewWriter(&buffer)
		_, err := writer.Write([]byte(s))
		if err != nil {
			return nil, err
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}
	return nil, fmt.Errorf("unsupported encoding %s", encoding)
}
