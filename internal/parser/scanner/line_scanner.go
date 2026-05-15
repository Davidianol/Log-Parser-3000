package scanner

import (
	"bufio"
	"io"
)

type LineScanner struct {
	scanner *bufio.Scanner
	lineNum int
}

func NewLineScanner(r io.Reader) *LineScanner {
	s := bufio.NewScanner(r)
	buf := make([]byte, 1024*1024)
	s.Buffer(buf, 1024*1024)
	return &LineScanner{scanner: s}
}

func (ls *LineScanner) Scan() bool {
	ok := ls.scanner.Scan()
	if ok {
		ls.lineNum++
	}
	return ok
}

func (ls *LineScanner) Text() string {
	return ls.scanner.Text()
}

func (ls *LineScanner) Err() error {
	return ls.scanner.Err()
}

func (ls *LineScanner) Line() int {
	return ls.lineNum
}
