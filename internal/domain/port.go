package domain

//START_PORTS
//NodeGuid,PortGuid,PortNum,MKey,GIDPrfx,MSMLID,LID,CapMsk,M_KeyLeasePeriod,DiagCode,LinkWidthActv,LinkWidthSup,LinkWidthEn,LocalPortNum,LinkSpeedEn,LinkSpeedActv,LMC,MKeyProtBits,LinkDownDefState,PortPhyState,PortState,LinkSpeedSup,VLArbHighCap,VLHighLimit,InitType,VLCap,MSMSL,NMTU,FilterRawOutb,FilterRawInb,PartEnfOutb,PartEnfInb,OpVLs,HoQLife,VLStallCnt,MTUCap,InitTypeReply,VLArbLowCap,PKeyViolations,MKeyViolations,SubnTmo,MulticastPKeyTrapSuppressionEnabled,ClientReregister,GUIDCap,QKeyViolations,MaxCreditHint,OverrunErrs,LocalPhyError,RespTimeValue,LinkRoundTripLatency,OOOSLMask,CapMsk2,FECActv,RetransActv
//0xhost1,0xhost1,1,0xffffffffff,0xffffffffff,1,1,2807162954,60,0,2,19,19,1,3841,2048,0,0,2,5,4,3841,8,0,0,3,0,5,0,0,0,0,2,0,0,5,0,8,0,0,18,1,0,8,0,0,8,8,16,270,0xffff,1074,14,0

type Port struct {
	ID     int `json:"id"`
	NodeID int `json:"node_id"`

	NodeGUID string `json:"node_guid"`
	PortGUID string `json:"port_guid"`
	PortNum  int    `json:"port_num"`

	MKey            string `json:"m_key"`
	GIDPrfx         string `json:"gid_prfx"`
	MSMLID          int    `json:"msmlid"`
	LID             int    `json:"lid"`
	CapMsk          uint64 `json:"cap_msk"`
	MKeyLeasePeriod int    `json:"m_key_lease_period"`
	DiagCode        int    `json:"diag_code"`

	LinkWidthActv int `json:"link_width_actv"`
	LinkWidthSup  int `json:"link_width_sup"`
	LinkWidthEn   int `json:"link_width_en"`
	LocalPortNum  int `json:"local_port_num"`

	LinkSpeedEn   int `json:"link_speed_en"`
	LinkSpeedActv int `json:"link_speed_actv"`
	LMC           int `json:"lmc"`
	MKeyProtBits  int `json:"m_key_prot_bits"`

	LinkDownDefState int `json:"link_down_def_state"`
	PortPhyState     int `json:"port_phy_state"`
	PortState        int `json:"port_state"`
	LinkSpeedSup     int `json:"link_speed_sup"`

	VLArbHighCap int `json:"vl_arb_high_cap"`
	VLHighLimit  int `json:"vl_high_limit"`
	InitType     int `json:"init_type"`
	VLCap        int `json:"vl_cap"`
	MSMSL        int `json:"msmsl"`
	NMTU         int `json:"nmtu"`

	FilterRawOutb int `json:"filter_raw_outb"`
	FilterRawInb  int `json:"filter_raw_inb"`
	PartEnfOutb   int `json:"part_enf_outb"`
	PartEnfInb    int `json:"part_enf_inb"`

	OpVLs      int `json:"op_vls"`
	HoQLife    int `json:"ho_q_life"`
	VLStallCnt int `json:"vl_stall_cnt"`
	MTUCap     int `json:"mtu_cap"`

	InitTypeReply int `json:"init_type_reply"`
	VLArbLowCap   int `json:"vl_arb_low_cap"`

	PKeyViolations int `json:"p_key_violations"`
	MKeyViolations int `json:"m_key_violations"`
	SubnTmo        int `json:"subn_tmo"`

	MulticastPKeyTrapSuppressionEnabled int `json:"multicast_p_key_trap_suppression_enabled"`
	ClientReregister                    int `json:"client_reregister"`
	GUIDCap                             int `json:"guid_cap"`
	QKeyViolations                      int `json:"q_key_violations"`
	MaxCreditHint                       int `json:"max_credit_hint"`
	OverrunErrs                         int `json:"overrun_errs"`
	LocalPhyError                       int `json:"local_phy_error"`
	RespTimeValue                       int `json:"resp_time_value"`

	LinkRoundTripLatency string `json:"link_round_trip_latency,omitempty"`
	OOOSLMask            string `json:"ooo_sl_mask,omitempty"`
	CapMsk2              string `json:"cap_msk2,omitempty"`
	FECActv              string `json:"fec_actv,omitempty"`
	RetransActv          string `json:"retrans_actv,omitempty"`
}
