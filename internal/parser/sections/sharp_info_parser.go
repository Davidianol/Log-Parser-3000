package sections

import (
	"fmt"
	"log_parser3000/internal/parser/scanner"
	"strings"

	"log_parser3000/internal/parser/raw"
)

func ParseSharpInfoSection(ls *scanner.LineScanner) ([]raw.SharpInfo, error) {
	var result []raw.SharpInfo
	var current *raw.SharpInfo

	for ls.Scan() {
		line := strings.TrimSpace(ls.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "START_") {
			return nil, fmt.Errorf("line %d: unexpected section start while parsing sharp info: %s", ls.Line(), line)
		}
		if strings.HasPrefix(line, "SW_GUID=") {
			if current != nil {
				result = append(result, *current)
			}
			nodeGUID := strings.TrimPrefix(line, "SW_GUID=")
			nodeGUID = strings.TrimSpace(nodeGUID)
			if !strings.HasPrefix(strings.ToLower(nodeGUID), "0x") {
				nodeGUID = "0x" + nodeGUID
			}
			current = &raw.SharpInfo{NodeGUID: nodeGUID}
			continue
		}
		if strings.HasPrefix(line, "----") {
			continue
		}
		if current == nil {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			switch key {
			case "endianness":
				current.Endianness = val
			case "enable_endianness_per_job":
				current.EnableEndiannessPerJob = val
			case "reproducibility_disable":
				current.ReproducibilityDisable = val
			}
		}
	}

	if current != nil {
		result = append(result, *current)
	}

	if err := ls.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
