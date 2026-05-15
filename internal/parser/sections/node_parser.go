package sections

import (
	"fmt"
	"log_parser3000/internal/parser/scanner"
	"strings"

	"log_parser3000/internal/parser/raw"
)

const nodesHeader = "NodeDesc,NumPorts,NodeType,ClassVersion,BaseVersion,SystemImageGUID,NodeGUID,PortGUID"
const countNodesHeader = 8

func ParseNodesSection(ls *scanner.LineScanner) ([]raw.Node, error) {
	if !ls.Scan() {
		return nil, fmt.Errorf("line %d: expected nodes header", ls.Line())
	}
	header := strings.TrimSpace(ls.Text())

	delimiter, err := detectDelimiter(header)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", ls.Line(), err)
	}
	if err := validateHeader(header, nodesHeader, delimiter); err != nil {
		return nil, fmt.Errorf("line %d: %w", ls.Line(), err)
	}

	var result []raw.Node

	for ls.Scan() {
		line := strings.TrimSpace(ls.Text())
		if line == "END_NODES" {
			return result, nil
		}
		if line == "" {
			continue
		}

		rec, err := parseCSVLine(line, delimiter)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse node: %w", ls.Line(), err)
		}
		if len(rec) != countNodesHeader {
			return nil, fmt.Errorf("line %d: invalid node field count: got=%d want=8", ls.Line(), len(rec))
		}

		result = append(result, raw.Node{
			NodeDesc:        rec[0],
			NumPorts:        rec[1],
			NodeType:        rec[2],
			ClassVersion:    rec[3],
			BaseVersion:     rec[4],
			SystemImageGUID: rec[5],
			NodeGUID:        rec[6],
			PortGUID:        rec[7],
		})
	}

	if err := ls.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("line %d: END_NODES not found", ls.Line())
}
