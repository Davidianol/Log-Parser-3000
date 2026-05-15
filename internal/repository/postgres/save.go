package postgres

import (
	"database/sql"
	"fmt"
	"log_parser3000/internal/domain"
	"time"
)

func (r *Repo) SaveParsedLog(parsed *domain.ParsedLog) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// 1. INSERT logs
	var logID int
	err = tx.QueryRow(`
		INSERT INTO logs (filename, status, node_count, uploaded_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		parsed.Log.Filename,
		string(domain.LogStatusDone),
		parsed.Log.NodeCount,
		time.Now(),
	).Scan(&logID)
	if err != nil {
		return 0, fmt.Errorf("insert log: %w", err)
	}

	// 2. INSERT nodes
	nodeIDByGUID := make(map[string]int, len(parsed.Nodes))
	for _, n := range parsed.Nodes {
		var nodeID int
		err = tx.QueryRow(`
			INSERT INTO nodes (log_id, node_desc, num_ports, node_type, system_image_guid, node_guid, port_guid)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`,
			logID, n.NodeDesc, n.NumPorts, int(n.NodeType),
			n.SystemImageGUID, n.NodeGUID, n.PortGUID,
		).Scan(&nodeID)
		if err != nil {
			return 0, fmt.Errorf("insert node guid=%s: %w", n.NodeGUID, err)
		}
		nodeIDByGUID[n.NodeGUID] = nodeID
	}

	// 3. INSERT ports
	for _, p := range parsed.Ports {
		nodeID, ok := nodeIDByGUID[p.NodeGUID]
		if !ok {
			return 0, fmt.Errorf("port references unknown node guid=%s", p.NodeGUID)
		}
		// А вы знаете, что такое безумие?
		_, err = tx.Exec(`
			INSERT INTO ports (
				node_id, log_id, node_guid, port_guid, port_num,
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
				link_round_trip_latency, ooo_sl_mask, cap_msk2, fec_actv, retrans_actv
			) VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,
				$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,
				$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,
				$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,
				$46,$47,$48,$49,$50,$51,$52,$53,$54,$55,$56
			)`,
			nodeID, logID, p.NodeGUID, p.PortGUID, p.PortNum,
			p.MKey, p.GIDPrfx, p.MSMLID, p.LID, p.CapMsk, p.MKeyLeasePeriod, p.DiagCode,
			p.LinkWidthActv, p.LinkWidthSup, p.LinkWidthEn, p.LocalPortNum,
			p.LinkSpeedEn, p.LinkSpeedActv, p.LMC, p.MKeyProtBits,
			p.LinkDownDefState, p.PortPhyState, p.PortState, p.LinkSpeedSup,
			p.VLArbHighCap, p.VLHighLimit, p.InitType, p.VLCap, p.MSMSL, p.NMTU,
			p.FilterRawOutb, p.FilterRawInb, p.PartEnfOutb, p.PartEnfInb,
			p.OpVLs, p.HoQLife, p.VLStallCnt, p.MTUCap,
			p.InitTypeReply, p.VLArbLowCap,
			p.PKeyViolations, p.MKeyViolations, p.SubnTmo,
			p.MulticastPKeyTrapSuppressionEnabled, p.ClientReregister,
			p.GUIDCap, p.QKeyViolations, p.MaxCreditHint,
			p.OverrunErrs, p.LocalPhyError, p.RespTimeValue,
			nullableString(p.LinkRoundTripLatency), nullableString(p.OOOSLMask),
			nullableString(p.CapMsk2), nullableString(p.FECActv), nullableString(p.RetransActv),
		)
		if err != nil {
			return 0, fmt.Errorf("insert port node_guid=%s port_num=%d: %w", p.NodeGUID, p.PortNum, err)
		}
	}

	// 4. INSERT nodes_info
	for _, ni := range parsed.NodesInfo {
		nodeGUID := nodeInfoGUID(&ni)
		if nodeGUID == "" {
			continue
		}
		nodeID, ok := nodeIDByGUID[nodeGUID]
		if !ok {
			continue
		}

		// GeneralInfo поля
		var giSerial, giPart, giRevision, giProduct sql.NullString
		if ni.GeneralInfo != nil {
			giSerial = nullableString(ni.GeneralInfo.SerialNumber)
			giPart = nullableString(ni.GeneralInfo.PartNumber)
			giRevision = nullableString(ni.GeneralInfo.Revision)
			giProduct = nullableString(ni.GeneralInfo.ProductName)
		}

		// SwitchInfo поля
		var (
			swLinearFDBCap, swRandomFDBCap, swMCastFDBCap, swLinearFDBTop sql.NullInt64
			swDefPort, swDefMCastPri, swDefMCastNotPri                    sql.NullInt64
			swLifeTime, swPortStateChange, swSLVL, swLidsPerPort          sql.NullInt64
			swPartEnf, swInbEnf, swOutbEnf                                sql.NullInt64
			swFilterRawInb, swFilterRawOutb, swENP0, swMCastFDBTop        sql.NullInt64
		)
		if ni.SwitchInfo != nil {
			swLinearFDBCap = nullableIntPtr(ni.SwitchInfo.LinearFDBCap)
			swRandomFDBCap = nullableIntPtr(ni.SwitchInfo.RandomFDBCap)
			swMCastFDBCap = nullableIntPtr(ni.SwitchInfo.MCastFDBCap)
			swLinearFDBTop = nullableIntPtr(ni.SwitchInfo.LinearFDBTop)
			swDefPort = nullableIntPtr(ni.SwitchInfo.DefPort)
			swDefMCastPri = nullableIntPtr(ni.SwitchInfo.DefMCastPriPort)
			swDefMCastNotPri = nullableIntPtr(ni.SwitchInfo.DefMCastNotPriPort)
			swLifeTime = nullableIntPtr(ni.SwitchInfo.LifeTimeValue)
			swPortStateChange = nullableIntPtr(ni.SwitchInfo.PortStateChange)
			swSLVL = nullableIntPtr(ni.SwitchInfo.OptimizedSLVLMapping)
			swLidsPerPort = nullableIntPtr(ni.SwitchInfo.LidsPerPort)
			swPartEnf = nullableIntPtr(ni.SwitchInfo.PartEnfCap)
			swInbEnf = nullableIntPtr(ni.SwitchInfo.InbEnfCap)
			swOutbEnf = nullableIntPtr(ni.SwitchInfo.OutbEnfCap)
			swFilterRawInb = nullableIntPtr(ni.SwitchInfo.FilterRawInbCap)
			swFilterRawOutb = nullableIntPtr(ni.SwitchInfo.FilterRawOutbCap)
			swENP0 = nullableIntPtr(ni.SwitchInfo.ENP0)
			swMCastFDBTop = nullableIntPtr(ni.SwitchInfo.MCastFDBTop)
		}

		// SharpInfo поля
		var sharpEndianness, sharpEnable, sharpRepro sql.NullInt64
		if ni.SharpInfo != nil {
			sharpEndianness = nullableIntPtr(ni.SharpInfo.Endianness)
			sharpEnable = nullableIntPtr(ni.SharpInfo.EnableEndiannessPerJob)
			sharpRepro = nullableIntPtr(ni.SharpInfo.ReproducibilityDisable)
		}

		_, err = tx.Exec(`
			INSERT INTO nodes_info (
				node_id, log_id,
				gi_serial_number, gi_part_number, gi_revision, gi_product_name,
				sw_linear_fdb_cap, sw_random_fdb_cap, sw_mcast_fdb_cap, sw_linear_fdb_top,
				sw_def_port, sw_def_mcast_pri_port, sw_def_mcast_not_pri_port,
				sw_life_time_value, sw_port_state_change, sw_optimized_slvl_mapping,
				sw_lids_per_port, sw_part_enf_cap, sw_inb_enf_cap, sw_outb_enf_cap,
				sw_filter_raw_inb_cap, sw_filter_raw_outb_cap, sw_enp0, sw_mcast_fdb_top,
				sharp_endianness, sharp_enable_endianness_per_job, sharp_reproducibility_disable
			) VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,
				$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27
			)`,
			nodeID, logID,
			giSerial, giPart, giRevision, giProduct,
			swLinearFDBCap, swRandomFDBCap, swMCastFDBCap, swLinearFDBTop,
			swDefPort, swDefMCastPri, swDefMCastNotPri,
			swLifeTime, swPortStateChange, swSLVL,
			swLidsPerPort, swPartEnf, swInbEnf, swOutbEnf,
			swFilterRawInb, swFilterRawOutb, swENP0, swMCastFDBTop,
			sharpEndianness, sharpEnable, sharpRepro,
		)
		if err != nil {
			return 0, fmt.Errorf("insert node_info node_guid=%s: %w", nodeGUID, err)
		}
	}

	// 5. INSERT topology_groups + topology_group_nodes
	for _, tg := range parsed.Topology {
		var groupID int
		err = tx.QueryRow(`
			INSERT INTO topology_groups (log_id, key) VALUES ($1, $2) RETURNING id`,
			logID, tg.Key,
		).Scan(&groupID)
		if err != nil {
			return 0, fmt.Errorf("insert topology group key=%s: %w", tg.Key, err)
		}
		for _, guid := range tg.NodeGUIDs {
			_, err = tx.Exec(`
				INSERT INTO topology_group_nodes (group_id, node_guid) VALUES ($1, $2)`,
				groupID, guid,
			)
			if err != nil {
				return 0, fmt.Errorf("insert topology node guid=%s: %w", guid, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}
	return logID, nil
}

func (r *Repo) SaveLogError(filename string, errMsg string) (int, error) {
	var logID int
	err := r.db.QueryRow(`
		INSERT INTO logs (filename, status, node_count, uploaded_at, error_message)
		VALUES ($1, $2, 0, NOW(), $3)
		RETURNING id`,
		filename,
		string(domain.LogStatusError),
		errMsg,
	).Scan(&logID)
	if err != nil {
		return 0, fmt.Errorf("insert log error: %w", err)
	}
	return logID, nil
}
