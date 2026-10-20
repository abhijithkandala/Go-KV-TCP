package resp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Reader struct {
	br *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{br: bufio.NewReader(rd)}
}

func (r *Reader) readLine() (string, error) {
	str, err := r.br.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(str, "\r\n"), nil
}

func (r *Reader) readBulkString() (string, error) {
	str, err := r.readLine()
	if err != nil {
		return "", err
	}
	if len(str) == 0 || str[0] != '$' {
		return "", errors.New("expected bulk string")
	}
	n, err := strconv.Atoi(str[1:])
	if err != nil {
		return "", errors.New("invalid bulk string length")
	}
	if n < 0 {
		return "", errors.New("null bulk string")
	}

	buf := make([]byte, n)
	if _, err := io.ReadFull(r.br, buf); err != nil {
		return "", err
	}
	if _, err := r.readLine(); err != nil {
		return "", err
	}
	return string(buf), nil
}

const maxCommandElements = 1024

func (r *Reader) ReadCommand() ([]string, error) {
	str, err := r.readLine()
	if err != nil {
		return nil, err
	}
	if len(str) == 0 || str[0] != '*' {
		return nil, errors.New("expected array")
	}

	n, err := strconv.Atoi(str[1:])
	if err != nil {
		return nil, errors.New("invalid array length")
	}
	if n <= 0 || n > maxCommandElements {
		return nil, errors.New("invalid array length")
	}

	parts := make([]string, n)
	for i := 0; i < n; i++ {
		s, err := r.readBulkString()
		if err != nil {
			return nil, err
		}
		parts[i] = s
	}
	return parts, nil
}
