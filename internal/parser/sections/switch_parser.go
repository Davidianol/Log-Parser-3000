package sections

import (
	"fmt"
	"log_parser3000/internal/parser/raw"
	"log_parser3000/internal/parser/scanner"
	"strings"
)

const switchesHeader = "NodeGUID,LinearFDBCap,RandomFDBCap,MCastFDBCap,LinearFDBTop,DefPort,DefMCastPriPort,DefMCastNotPriPort,LifeTimeValue,PortStateChange,OptimizedSLVLMapping,LidsPerPort,PartEnfCap,InbEnfCap,OutbEnfCap,FilterRawInbCap,FilterRawOutbCap,ENP0,MCastFDBTop"
const countSwitchesHeader = 19

func ParseSwitchesSection(ls *scanner.LineScanner) ([]raw.SwitchInfo, error) {
	if !ls.Scan() {
		return nil, fmt.Errorf("line %d: expected switches header", ls.Line())
	}
	header := strings.TrimSpace(ls.Text())

	delimiter, err := detectDelimiter(header)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", ls.Line(), err)
	}
	if err := validateHeader(header, switchesHeader, delimiter); err != nil {
		return nil, fmt.Errorf("line %d: %w", ls.Line(), err)
	}

	var result []raw.SwitchInfo

	for ls.Scan() {
		line := strings.TrimSpace(ls.Text())
		if line == "END_SWITCHES" {
			return result, nil
		}
		if line == "" {
			continue
		}

		rec, err := parseCSVLine(line, delimiter)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse switch info: %w", ls.Line(), err)
		}
		if len(rec) != countSwitchesHeader {
			return nil, fmt.Errorf("line %d: invalid switch info field count: got=%d want=19", ls.Line(), len(rec))
		}

		result = append(result, raw.SwitchInfo{
			NodeGUID:             rec[0],
			LinearFDBCap:         rec[1],
			RandomFDBCap:         rec[2],
			MCastFDBCap:          rec[3],
			LinearFDBTop:         rec[4],
			DefPort:              rec[5],
			DefMCastPriPort:      rec[6],
			DefMCastNotPriPort:   rec[7],
			LifeTimeValue:        rec[8],
			PortStateChange:      rec[9],
			OptimizedSLVLMapping: rec[10],
			LidsPerPort:          rec[11],
			PartEnfCap:           rec[12],
			InbEnfCap:            rec[13],
			OutbEnfCap:           rec[14],
			FilterRawInbCap:      rec[15],
			FilterRawOutbCap:     rec[16],
			ENP0:                 rec[17],
			MCastFDBTop:          rec[18],
		})
	}

	if err := ls.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("line %d: END_SWITCHES not found", ls.Line())
}
