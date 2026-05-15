package parser

import (
	"fmt"
	"log_parser3000/internal/domain"
	"log_parser3000/internal/parser/raw"
	"sort"
	"strconv"
	"strings"
)

func MapParsedDataToDomain(filename string, db *ParsedDBCSV, sharp []raw.SharpInfo) (*domain.ParsedLog, error) {
	nodes := make([]domain.Node, 0, len(db.Nodes))
	ports := make([]domain.Port, 0, len(db.Ports))

	switchByGUID := make(map[string]raw.SwitchInfo, len(db.Switches))
	for _, sw := range db.Switches {
		switchByGUID[normalizeGUID(sw.NodeGUID)] = sw
	}

	generalByGUID := make(map[string]raw.GeneralInfo, len(db.GeneralInfo))
	for _, gi := range db.GeneralInfo {
		generalByGUID[normalizeGUID(gi.NodeGUID)] = gi
	}

	sharpByGUID := make(map[string]raw.SharpInfo, len(sharp))
	for _, si := range sharp {
		sharpByGUID[normalizeGUID(si.NodeGUID)] = si
	}

	nodeGUIDSet := make(map[string]struct{}, len(db.Nodes))

	for _, rn := range db.Nodes {
		n, err := mapNode(rn)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
		nodeGUIDSet[n.NodeGUID] = struct{}{}
	}

	for _, rp := range db.Ports {
		p, err := mapPort(rp)
		if err != nil {
			return nil, err
		}
		if _, ok := nodeGUIDSet[p.NodeGUID]; !ok {
			return nil, fmt.Errorf("port references unknown node guid=%s", p.NodeGUID)
		}
		ports = append(ports, p)
	}

	nodesInfo := buildNodeInfos(nodes, switchByGUID, generalByGUID, sharpByGUID)
	topology := buildTopology(nodes, ports)

	return &domain.ParsedLog{
		Log: domain.Log{
			Filename:  filename,
			Status:    domain.LogStatusDone,
			NodeCount: len(nodes),
		},
		Nodes:     nodes,
		Ports:     ports,
		NodesInfo: nodesInfo,
		Topology:  topology,
	}, nil
}

func mapNode(r raw.Node) (domain.Node, error) {
	numPorts, err := parseIntStrict("NumPorts", r.NumPorts, r.NodeGUID)
	if err != nil {
		return domain.Node{}, err
	}
	nodeType, err := parseIntStrict("NodeType", r.NodeType, r.NodeGUID)
	if err != nil {
		return domain.Node{}, err
	}

	return domain.Node{
		NodeDesc:        strings.TrimSpace(r.NodeDesc),
		NumPorts:        numPorts,
		NodeType:        domain.NodeType(nodeType),
		SystemImageGUID: normalizeGUID(r.SystemImageGUID),
		NodeGUID:        normalizeGUID(r.NodeGUID),
		PortGUID:        normalizeGUID(r.PortGUID),
	}, nil
}

func mapPort(r raw.Port) (domain.Port, error) {
	nodeGUID := normalizeGUID(r.NodeGUID)

	portNum, err := parseIntStrict("PortNum", r.PortNum, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	msmlid, err := parseIntStrict("MSMLID", r.MSMLID, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	lid, err := parseIntStrict("LID", r.LID, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	capMask, err := parseUint64Flexible("CapMsk", r.CapMsk, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	mKeyLeasePeriod, err := parseIntStrict("M_KeyLeasePeriod", r.MKeyLeasePeriod, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	diagCode, err := parseIntStrict("DiagCode", r.DiagCode, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	linkWidthActv, err := parseIntStrict("LinkWidthActv", r.LinkWidthActv, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	linkWidthSup, err := parseIntStrict("LinkWidthSup", r.LinkWidthSup, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	linkWidthEn, err := parseIntStrict("LinkWidthEn", r.LinkWidthEn, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	localPortNum, err := parseIntStrict("LocalPortNum", r.LocalPortNum, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	linkSpeedEn, err := parseIntStrict("LinkSpeedEn", r.LinkSpeedEn, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	linkSpeedActv, err := parseIntStrict("LinkSpeedActv", r.LinkSpeedActv, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	lmc, err := parseIntStrict("LMC", r.LMC, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	mKeyProtBits, err := parseIntStrict("MKeyProtBits", r.MKeyProtBits, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	linkDownDefState, err := parseIntStrict("LinkDownDefState", r.LinkDownDefState, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	portPhyState, err := parseIntStrict("PortPhyState", r.PortPhyState, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	portState, err := parseIntStrict("PortState", r.PortState, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	linkSpeedSup, err := parseIntStrict("LinkSpeedSup", r.LinkSpeedSup, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	vlArbHighCap, err := parseIntStrict("VLArbHighCap", r.VLArbHighCap, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	vlHighLimit, err := parseIntStrict("VLHighLimit", r.VLHighLimit, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	initType, err := parseIntStrict("InitType", r.InitType, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	vlCap, err := parseIntStrict("VLCap", r.VLCap, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	msmsl, err := parseIntStrict("MSMSL", r.MSMSL, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	nmtu, err := parseIntStrict("NMTU", r.NMTU, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	filterRawOutb, err := parseIntStrict("FilterRawOutb", r.FilterRawOutb, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	filterRawInb, err := parseIntStrict("FilterRawInb", r.FilterRawInb, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	partEnfOutb, err := parseIntStrict("PartEnfOutb", r.PartEnfOutb, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	partEnfInb, err := parseIntStrict("PartEnfInb", r.PartEnfInb, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	opVLs, err := parseIntStrict("OpVLs", r.OpVLs, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	hoQLife, err := parseIntStrict("HoQLife", r.HoQLife, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	vlStallCnt, err := parseIntStrict("VLStallCnt", r.VLStallCnt, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	mtuCap, err := parseIntStrict("MTUCap", r.MTUCap, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	initTypeReply, err := parseIntStrict("InitTypeReply", r.InitTypeReply, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	vlArbLowCap, err := parseIntStrict("VLArbLowCap", r.VLArbLowCap, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	pKeyViolations, err := parseIntStrict("PKeyViolations", r.PKeyViolations, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	mKeyViolations, err := parseIntStrict("MKeyViolations", r.MKeyViolations, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	subnTmo, err := parseIntStrict("SubnTmo", r.SubnTmo, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	multiTrap, err := parseIntStrict("MulticastPKeyTrapSuppressionEnabled", r.MulticastPKeyTrapSuppressionEnabled, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	clientReregister, err := parseIntStrict("ClientReregister", r.ClientReregister, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	guidCap, err := parseIntStrict("GUIDCap", r.GUIDCap, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	qKeyViolations, err := parseIntStrict("QKeyViolations", r.QKeyViolations, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	maxCreditHint, err := parseIntStrict("MaxCreditHint", r.MaxCreditHint, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	overrunErrs, err := parseIntStrict("OverrunErrs", r.OverrunErrs, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	localPhyError, err := parseIntStrict("LocalPhyError", r.LocalPhyError, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}
	respTimeValue, err := parseIntStrict("RespTimeValue", r.RespTimeValue, nodeGUID)
	if err != nil {
		return domain.Port{}, err
	}

	return domain.Port{
		NodeGUID: nodeGUID,
		PortGUID: normalizeGUID(r.PortGUID),
		PortNum:  portNum,

		MKey:            strings.TrimSpace(r.MKey),
		GIDPrfx:         strings.TrimSpace(r.GIDPrfx),
		MSMLID:          msmlid,
		LID:             lid,
		CapMsk:          capMask,
		MKeyLeasePeriod: mKeyLeasePeriod,
		DiagCode:        diagCode,

		LinkWidthActv: linkWidthActv,
		LinkWidthSup:  linkWidthSup,
		LinkWidthEn:   linkWidthEn,
		LocalPortNum:  localPortNum,

		LinkSpeedEn:   linkSpeedEn,
		LinkSpeedActv: linkSpeedActv,
		LMC:           lmc,
		MKeyProtBits:  mKeyProtBits,

		LinkDownDefState: linkDownDefState,
		PortPhyState:     portPhyState,
		PortState:        portState,
		LinkSpeedSup:     linkSpeedSup,

		VLArbHighCap: vlArbHighCap,
		VLHighLimit:  vlHighLimit,
		InitType:     initType,
		VLCap:        vlCap,
		MSMSL:        msmsl,
		NMTU:         nmtu,

		FilterRawOutb: filterRawOutb,
		FilterRawInb:  filterRawInb,
		PartEnfOutb:   partEnfOutb,
		PartEnfInb:    partEnfInb,

		OpVLs:      opVLs,
		HoQLife:    hoQLife,
		VLStallCnt: vlStallCnt,
		MTUCap:     mtuCap,

		InitTypeReply: initTypeReply,
		VLArbLowCap:   vlArbLowCap,

		PKeyViolations: pKeyViolations,
		MKeyViolations: mKeyViolations,
		SubnTmo:        subnTmo,

		MulticastPKeyTrapSuppressionEnabled: multiTrap,
		ClientReregister:                    clientReregister,
		GUIDCap:                             guidCap,
		QKeyViolations:                      qKeyViolations,
		MaxCreditHint:                       maxCreditHint,
		OverrunErrs:                         overrunErrs,
		LocalPhyError:                       localPhyError,
		RespTimeValue:                       respTimeValue,

		LinkRoundTripLatency: nullableStringValue(r.LinkRoundTripLatency),
		OOOSLMask:            nullableStringValue(r.OOOSLMask),
		CapMsk2:              nullableStringValue(r.CapMsk2),
		FECActv:              nullableStringValue(r.FECActv),
		RetransActv:          nullableStringValue(r.RetransActv),
	}, nil
}

func buildNodeInfos(
	nodes []domain.Node,
	switchByGUID map[string]raw.SwitchInfo,
	generalByGUID map[string]raw.GeneralInfo,
	sharpByGUID map[string]raw.SharpInfo,
) []domain.NodeInfo {
	out := make([]domain.NodeInfo, 0, len(nodes))

	for _, n := range nodes {
		info := domain.NodeInfo{
			NodeGUID: n.NodeGUID,
		}

		if g, ok := generalByGUID[n.NodeGUID]; ok {
			info.GeneralInfo = &domain.GeneralInfo{
				NodeGUID:     normalizeGUID(g.NodeGUID),
				SerialNumber: strings.TrimSpace(g.SerialNumber),
				PartNumber:   strings.TrimSpace(g.PartNumber),
				Revision:     strings.TrimSpace(g.Revision),
				ProductName:  strings.TrimSpace(g.ProductName),
			}
		}

		if sw, ok := switchByGUID[n.NodeGUID]; ok {
			info.SwitchInfo = &domain.SwitchInfo{
				NodeGUID:             normalizeGUID(sw.NodeGUID),
				LinearFDBCap:         nullableIntValue(sw.LinearFDBCap),
				RandomFDBCap:         nullableIntValue(sw.RandomFDBCap),
				MCastFDBCap:          nullableIntValue(sw.MCastFDBCap),
				LinearFDBTop:         nullableIntValue(sw.LinearFDBTop),
				DefPort:              nullableIntValue(sw.DefPort),
				DefMCastPriPort:      nullableIntValue(sw.DefMCastPriPort),
				DefMCastNotPriPort:   nullableIntValue(sw.DefMCastNotPriPort),
				LifeTimeValue:        nullableIntValue(sw.LifeTimeValue),
				PortStateChange:      nullableIntValue(sw.PortStateChange),
				OptimizedSLVLMapping: nullableIntValue(sw.OptimizedSLVLMapping),
				LidsPerPort:          nullableIntValue(sw.LidsPerPort),
				PartEnfCap:           nullableIntValue(sw.PartEnfCap),
				InbEnfCap:            nullableIntValue(sw.InbEnfCap),
				OutbEnfCap:           nullableIntValue(sw.OutbEnfCap),
				FilterRawInbCap:      nullableIntValue(sw.FilterRawInbCap),
				FilterRawOutbCap:     nullableIntValue(sw.FilterRawOutbCap),
				ENP0:                 nullableIntValue(sw.ENP0),
				MCastFDBTop:          nullableIntValue(sw.MCastFDBTop),
			}
		}

		if sh, ok := sharpByGUID[n.NodeGUID]; ok {
			info.SharpInfo = &domain.SharpInfo{
				NodeGUID:               normalizeGUID(sh.NodeGUID),
				Endianness:             nullableIntValue(sh.Endianness),
				EnableEndiannessPerJob: nullableIntValue(sh.EnableEndiannessPerJob),
				ReproducibilityDisable: nullableIntValue(sh.ReproducibilityDisable),
			}
		}

		if info.SwitchInfo != nil || info.GeneralInfo != nil || info.SharpInfo != nil {
			out = append(out, info)
		}
	}

	return out
}

func buildTopology(nodes []domain.Node, ports []domain.Port) []domain.TopologyGroup {
	// Ищем узлы, у которых есть хотя бы один активный порт
	activeNodes := make(map[string]bool)
	for _, p := range ports {
		if p.PortState == 4 { // 4 = PortActive
			activeNodes[p.NodeGUID] = true
		}
	}

	groups := make(map[string][]string) // key -> []NodeGUID

	// Формируем группы по типу узлов и состояния портов
	for _, n := range nodes {
		baseKey := n.NodeType.String() // "host", "switch", "router"

		statusSuffix := "_isolated"
		if activeNodes[n.NodeGUID] {
			statusSuffix = "_active"
		}

		key := baseKey + statusSuffix
		groups[key] = append(groups[key], n.NodeGUID)
	}

	// Стабильная сортировка для детерминированного вывода API
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]domain.TopologyGroup, 0, len(groups))
	for _, k := range keys {
		guids := groups[k]
		sort.Strings(guids)
		result = append(result, domain.TopologyGroup{
			Key:       k,
			NodeGUIDs: guids,
		})
	}

	return result
}

func parseIntStrict(field, value, guid string) (int, error) {
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("%s for guid=%s: %w", field, normalizeGUID(guid), err)
	}
	return v, nil
}

func parseUint64Flexible(field, value, guid string) (uint64, error) {
	s := strings.TrimSpace(value)
	base := 10
	if strings.HasPrefix(strings.ToLower(s), "0x") {
		base = 16
		s = s[2:]
	}
	v, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		return 0, fmt.Errorf("%s for guid=%s: %w", field, normalizeGUID(guid), err)
	}
	return v, nil
}

func nullableIntValue(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, "N/A") {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}

func nullableStringValue(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, "N/A") {
		return ""
	}
	return s
}

func normalizeGUID(v string) string {
	s := strings.TrimSpace(v)
	if s == "" || strings.EqualFold(s, "N/A") {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(s), "0x") {
		return strings.ToLower(s)
	}
	return "0x" + strings.ToLower(s)
}
