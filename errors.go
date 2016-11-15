package zif

import (
	"bufio"
	"io"
)

type ErrorReader struct {
	reader *bufio.Reader
	err    error
}

func NewErrorReader(r io.Reader) *ErrorReader {
	return &ErrorReader{bufio.NewReader(r), nil}
}

func (er *ErrorReader) ReadString(delim byte) string {
	var ret string

	ret, er.err = er.reader.ReadString(delim)

	if er.err != nil {
		return ""
	}

	return ret
}

type AddressResolutionError struct {
	address string
}

func (a AddressResolutionError) Error() string {
	return "Failed to resolve address, address may not exist or is not reachable"
}
