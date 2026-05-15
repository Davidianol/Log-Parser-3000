package sections

import (
	"fmt"
	"log_parser3000/internal/parser/scanner"
	"strings"

	"log_parser3000/internal/parser/raw"
)

const generalInfoHeader = "NodeGuid,SerialNumber,PartNumber,Revision,ProductName"
const countGeneralInfo = 5

func ParseGeneralInfoSection(ls *scanner.LineScanner) ([]raw.GeneralInfo, error) {
	if !ls.Scan() {
		return nil, fmt.Errorf("line %d: expected general info header", ls.Line())
	}
	header := strings.TrimSpace(ls.Text())

	delimiter, err := detectDelimiter(header)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", ls.Line(), err)
	}
	if err := validateHeader(header, generalInfoHeader, delimiter); err != nil {
		return nil, fmt.Errorf("line %d: %w", ls.Line(), err)
	}

	var result []raw.GeneralInfo

	for ls.Scan() {
		line := strings.TrimSpace(ls.Text())
		if line == "END_SYSTEM_GENERAL_INFORMATION" {
			return result, nil
		}
		if line == "" {
			continue
		}

		rec, err := parseCSVLine(line, delimiter)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse general info: %w", ls.Line(), err)
		}
		if len(rec) != countGeneralInfo {
			return nil, fmt.Errorf("line %d: invalid general info field count: got=%d want=5", ls.Line(), len(rec))
		}

		result = append(result, raw.GeneralInfo{
			NodeGUID:     rec[0],
			SerialNumber: rec[1],
			PartNumber:   rec[2],
			Revision:     rec[3],
			ProductName:  rec[4],
		})
	}

	if err := ls.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("line %d: END_SYSTEM_GENERAL_INFORMATION not found", ls.Line())
}
