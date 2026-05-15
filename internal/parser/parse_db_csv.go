package parser

import (
	"fmt"
	"io"
	"log_parser3000/internal/parser/scanner"
	"log_parser3000/internal/parser/sections"
	"strings"

	"log_parser3000/internal/parser/raw"
)

type ParsedDBCSV struct {
	Nodes       []raw.Node
	Ports       []raw.Port
	Switches    []raw.SwitchInfo
	GeneralInfo []raw.GeneralInfo
}

func ParseDBCSV(r io.Reader) (*ParsedDBCSV, error) {
	ls := scanner.NewLineScanner(r)

	out := &ParsedDBCSV{}

	for ls.Scan() {
		line := strings.TrimSpace(ls.Text())
		if line == "" {
			continue
		}

		switch line {
		case "START_NODES":
			nodes, err := sections.ParseNodesSection(ls)
			if err != nil {
				return nil, err
			}
			out.Nodes = nodes

		case "START_PORTS":
			ports, err := sections.ParsePortsSection(ls)
			if err != nil {
				return nil, err
			}
			out.Ports = ports

		case "START_SWITCHES":
			switches, err := sections.ParseSwitchesSection(ls)
			if err != nil {
				return nil, err
			}
			out.Switches = switches

		case "START_SYSTEM_GENERAL_INFORMATION":
			info, err := sections.ParseGeneralInfoSection(ls)
			if err != nil {
				return nil, err
			}
			out.GeneralInfo = info
		}
	}

	if err := ls.Err(); err != nil {
		return nil, err
	}

	if len(out.Nodes) == 0 {
		return nil, fmt.Errorf("START_NODES section not found or empty")
	}
	if len(out.Ports) == 0 {
		return nil, fmt.Errorf("START_PORTS section not found or empty")
	}

	return out, nil
}
