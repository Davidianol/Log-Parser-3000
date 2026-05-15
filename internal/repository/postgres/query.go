package postgres

import (
	"database/sql"
	"fmt"
	"log_parser3000/internal/domain"
)

func (r *Repo) GetLogByID(id int) (*domain.Log, error) {
	var l domain.Log
	var errMsg sql.NullString

	err := r.db.QueryRow(`
		SELECT id, filename, status, node_count, uploaded_at, error_message
		FROM logs WHERE id = $1`, id,
	).Scan(&l.ID, &l.Filename, &l.Status, &l.NodeCount, &l.UploadedAt, &errMsg)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("log id=%d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("get log: %w", err)
	}
	if errMsg.Valid {
		l.ErrorMessage = &errMsg.String
	}
	return &l, nil
}

func (r *Repo) GetNodeByID(id int) (*domain.Node, error) {
	var n domain.Node
	err := r.db.QueryRow(`
		SELECT id, log_id, node_desc, num_ports, node_type, system_image_guid, node_guid, port_guid
		FROM nodes WHERE id = $1`, id,
	).Scan(&n.ID, &n.LogID, &n.NodeDesc, &n.NumPorts, &n.NodeType,
		&n.SystemImageGUID, &n.NodeGUID, &n.PortGUID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("node id=%d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("get node: %w", err)
	}
	return &n, nil
}

func (r *Repo) GetNodesByLogID(logID int) ([]domain.Node, error) {
	rows, err := r.db.Query(`
		SELECT id, log_id, node_desc, num_ports, node_type, system_image_guid, node_guid, port_guid
		FROM nodes WHERE log_id = $1
		ORDER BY id`, logID,
	)
	if err != nil {
		return nil, fmt.Errorf("get nodes: %w", err)
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var n domain.Node
		if err := rows.Scan(&n.ID, &n.LogID, &n.NodeDesc, &n.NumPorts, &n.NodeType,
			&n.SystemImageGUID, &n.NodeGUID, &n.PortGUID); err != nil {
			return nil, fmt.Errorf("scan node: %w", err)
		}
		nodes = append(nodes, n)
	}
	return nodes, rows.Err()
}

func (r *Repo) GetPortsByNodeID(nodeID int) ([]domain.Port, error) {
	// Проверяем существование узла
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM nodes WHERE id=$1)`, nodeID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("node id=%d not found", nodeID)
	}

	rows, err := r.db.Query(`
		SELECT id, node_id, node_guid, port_guid, port_num,
			m_key, gid_prfx, msmlid, lid, cap_msk, m_key_lease_period, diag_code,
			link_width_actv, link_width_sup, link_width_en, local_port_num,
			link_speed_en, link_speed_actv, lmc, m_key_prot_bits,
			link_down_def_state, port_phy_state, port_state, link_speed_sup,
			vl_arb_high_cap, vl_high_limit, init_type, vl_cap, msmsl, nmtu,
			filter_raw_outb, filter_raw_inb, part_enf_outb, part_enf_inb,
			op_vls, ho_q_life, vl_stall_cnt, mtu_cap,
			init_type_reply, vl_arb_low_cap,
			p_key_violations, m_key_violations, subn_tmo,
			multicast_p_key_trap_suppression_enabled, client_reregister,
			guid_cap, q_key_violations, max_credit_hint,
			overrun_errs, local_phy_error, resp_time_value,
			COALESCE(link_round_trip_latency, ''),
			COALESCE(ooo_sl_mask, ''),
			COALESCE(cap_msk2, ''),
			COALESCE(fec_actv, ''),
			COALESCE(retrans_actv, '')
		FROM ports WHERE node_id = $1
		ORDER BY port_num`, nodeID,
	)
	if err != nil {
		return nil, fmt.Errorf("get ports: %w", err)
	}
	defer rows.Close()

	ports := make([]domain.Port, 0)
	for rows.Next() {
		var p domain.Port
		if err := rows.Scan(
			&p.ID, &p.NodeID, &p.NodeGUID, &p.PortGUID, &p.PortNum,
			&p.MKey, &p.GIDPrfx, &p.MSMLID, &p.LID, &p.CapMsk, &p.MKeyLeasePeriod, &p.DiagCode,
			&p.LinkWidthActv, &p.LinkWidthSup, &p.LinkWidthEn, &p.LocalPortNum,
			&p.LinkSpeedEn, &p.LinkSpeedActv, &p.LMC, &p.MKeyProtBits,
			&p.LinkDownDefState, &p.PortPhyState, &p.PortState, &p.LinkSpeedSup,
			&p.VLArbHighCap, &p.VLHighLimit, &p.InitType, &p.VLCap, &p.MSMSL, &p.NMTU,
			&p.FilterRawOutb, &p.FilterRawInb, &p.PartEnfOutb, &p.PartEnfInb,
			&p.OpVLs, &p.HoQLife, &p.VLStallCnt, &p.MTUCap,
			&p.InitTypeReply, &p.VLArbLowCap,
			&p.PKeyViolations, &p.MKeyViolations, &p.SubnTmo,
			&p.MulticastPKeyTrapSuppressionEnabled, &p.ClientReregister,
			&p.GUIDCap, &p.QKeyViolations, &p.MaxCreditHint,
			&p.OverrunErrs, &p.LocalPhyError, &p.RespTimeValue,
			&p.LinkRoundTripLatency, &p.OOOSLMask, &p.CapMsk2, &p.FECActv, &p.RetransActv,
		); err != nil {
			return nil, fmt.Errorf("scan port: %w", err)
		}
		ports = append(ports, p)
	}
	return ports, rows.Err()
}

func (r *Repo) GetTopologyByLogID(logID int) ([]domain.TopologyGroup, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM logs WHERE id=$1)`, logID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("log id=%d not found", logID)
	}

	rows, err := r.db.Query(`
		SELECT tg.key, tgn.node_guid
		FROM topology_groups tg
		JOIN topology_group_nodes tgn ON tgn.group_id = tg.id
		WHERE tg.log_id = $1
		ORDER BY tg.key, tgn.node_guid`, logID,
	)
	if err != nil {
		return nil, fmt.Errorf("get topology: %w", err)
	}
	defer rows.Close()

	groupMap := make(map[string]*domain.TopologyGroup)
	var order []string

	for rows.Next() {
		var key, guid string
		if err := rows.Scan(&key, &guid); err != nil {
			return nil, fmt.Errorf("scan topology: %w", err)
		}
		if _, ok := groupMap[key]; !ok {
			groupMap[key] = &domain.TopologyGroup{Key: key}
			order = append(order, key)
		}
		groupMap[key].NodeGUIDs = append(groupMap[key].NodeGUIDs, guid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]domain.TopologyGroup, 0, len(order))
	for _, k := range order {
		result = append(result, *groupMap[k])
	}
	return result, nil
}
