package parser

import (
	"io"
	"log_parser3000/internal/parser/raw"
	"log_parser3000/internal/parser/scanner"
	"log_parser3000/internal/parser/sections"
)

func ParseSharpInfo(r io.Reader) ([]raw.SharpInfo, error) {
	ls := scanner.NewLineScanner(r)
	return sections.ParseSharpInfoSection(ls)
}
