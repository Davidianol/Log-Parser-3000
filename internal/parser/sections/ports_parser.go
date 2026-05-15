package sections

import (
	"fmt"
	"log_parser3000/internal/parser/scanner"
	"strings"

	"log_parser3000/internal/parser/raw"
)

const portsHeader = "NodeGuid,PortGuid,PortNum,MKey,GIDPrfx,MSMLID,LID,CapMsk,M_KeyLeasePeriod,DiagCode,LinkWidthActv,LinkWidthSup,LinkWidthEn,LocalPortNum,LinkSpeedEn,LinkSpeedActv,LMC,MKeyProtBits,LinkDownDefState,PortPhyState,PortState,LinkSpeedSup,VLArbHighCap,VLHighLimit,InitType,VLCap,MSMSL,NMTU,FilterRawOutb,FilterRawInb,PartEnfOutb,PartEnfInb,OpVLs,HoQLife,VLStallCnt,MTUCap,InitTypeReply,VLArbLowCap,PKeyViolations,MKeyViolations,SubnTmo,MulticastPKeyTrapSuppressionEnabled,ClientReregister,GUIDCap,QKeyViolations,MaxCreditHint,OverrunErrs,LocalPhyError,RespTimeValue,LinkRoundTripLatency,OOOSLMask,CapMsk2,FECActv,RetransActv"
const countPortsHeader = 54

func ParsePortsSection(ls *scanner.LineScanner) ([]raw.Port, error) {
	if !ls.Scan() {
		return nil, fmt.Errorf("line %d: expected ports header", ls.Line())
	}
	header := strings.TrimSpace(ls.Text())

	delimiter, err := detectDelimiter(header)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", ls.Line(), err)
	}
	if err := validateHeader(header, portsHeader, delimiter); err != nil {
		return nil, fmt.Errorf("line %d: %w", ls.Line(), err)
	}

	var result []raw.Port

	for ls.Scan() {
		line := strings.TrimSpace(ls.Text())
		if line == "END_PORTS" {
			return result, nil
		}
		if line == "" {
			continue
		}

		rec, err := parseCSVLine(line, delimiter)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse port: %w", ls.Line(), err)
		}
		if len(rec) != countPortsHeader {
			return nil, fmt.Errorf("line %d: invalid port field count: got=%d want=54", ls.Line(), len(rec))
		}

		result = append(result, raw.Port{
			NodeGUID:                            rec[0],
			PortGUID:                            rec[1],
			PortNum:                             rec[2],
			MKey:                                rec[3],
			GIDPrfx:                             rec[4],
			MSMLID:                              rec[5],
			LID:                                 rec[6],
			CapMsk:                              rec[7],
			MKeyLeasePeriod:                     rec[8],
			DiagCode:                            rec[9],
			LinkWidthActv:                       rec[10],
			LinkWidthSup:                        rec[11],
			LinkWidthEn:                         rec[12],
			LocalPortNum:                        rec[13],
			LinkSpeedEn:                         rec[14],
			LinkSpeedActv:                       rec[15],
			LMC:                                 rec[16],
			MKeyProtBits:                        rec[17],
			LinkDownDefState:                    rec[18],
			PortPhyState:                        rec[19],
			PortState:                           rec[20],
			LinkSpeedSup:                        rec[21],
			VLArbHighCap:                        rec[22],
			VLHighLimit:                         rec[23],
			InitType:                            rec[24],
			VLCap:                               rec[25],
			MSMSL:                               rec[26],
			NMTU:                                rec[27],
			FilterRawOutb:                       rec[28],
			FilterRawInb:                        rec[29],
			PartEnfOutb:                         rec[30],
			PartEnfInb:                          rec[31],
			OpVLs:                               rec[32],
			HoQLife:                             rec[33],
			VLStallCnt:                          rec[34],
			MTUCap:                              rec[35],
			InitTypeReply:                       rec[36],
			VLArbLowCap:                         rec[37],
			PKeyViolations:                      rec[38],
			MKeyViolations:                      rec[39],
			SubnTmo:                             rec[40],
			MulticastPKeyTrapSuppressionEnabled: rec[41],
			ClientReregister:                    rec[42],
			GUIDCap:                             rec[43],
			QKeyViolations:                      rec[44],
			MaxCreditHint:                       rec[45],
			OverrunErrs:                         rec[46],
			LocalPhyError:                       rec[47],
			RespTimeValue:                       rec[48],
			LinkRoundTripLatency:                rec[49],
			OOOSLMask:                           rec[50],
			CapMsk2:                             rec[51],
			FECActv:                             rec[52],
			RetransActv:                         rec[53],
		})
	}

	if err := ls.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("line %d: END_PORTS not found", ls.Line())
}
