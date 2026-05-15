package domain

type NodeInfo struct {
	ID       int    `json:"id"`
	NodeID   int    `json:"node_id"`
	NodeGUID string `json:"node_guid"`

	SwitchInfo  *SwitchInfo  `json:"switch_info,omitempty"`
	GeneralInfo *GeneralInfo `json:"general_info,omitempty"`
	SharpInfo   *SharpInfo   `json:"sharp_info,omitempty"`
}

//START_SYSTEM_GENERAL_INFORMATION
//NodeGuid,SerialNumber,PartNumber,Revision,ProductName
//0xswitch1,SOS123,MMM-MAV,AA,"Gorilla"

type GeneralInfo struct {
	NodeGUID     string `json:"node_guid"`
	SerialNumber string `json:"serial_number"`
	PartNumber   string `json:"part_number"`
	Revision     string `json:"revision"`
	ProductName  string `json:"product_name"`
}

//SW_GUID=switch3
//-------------------------------------------------------------------------------------------
//endianness = 0
//enable_endianness_per_job = 0
//reproducibility_disable = 4

type SharpInfo struct {
	NodeGUID               string `json:"node_guid"`
	Endianness             *int   `json:"endianness,omitempty"`
	EnableEndiannessPerJob *int   `json:"enable_endianness_per_job,omitempty"`
	ReproducibilityDisable *int   `json:"reproducibility_disable,omitempty"`
}

// START_SWITCHES
// NodeGUID,LinearFDBCap,RandomFDBCap,MCastFDBCap,LinearFDBTop,DefPort,DefMCastPriPort,DefMCastNotPriPort,LifeTimeValue,PortStateChange,OptimizedSLVLMapping,LidsPerPort,PartEnfCap,InbEnfCap,OutbEnfCap,FilterRawInbCap,FilterRawOutbCap,ENP0,MCastFDBTop
// 0xswitch1,49152,0,8192,78,0,255,255,18,0,3,0,32,1,1,1,1,0,49183
type SwitchInfo struct {
	NodeGUID             string `json:"node_guid"`
	LinearFDBCap         *int   `json:"linear_fdb_cap,omitempty"`
	RandomFDBCap         *int   `json:"random_fdb_cap,omitempty"`
	MCastFDBCap          *int   `json:"mcast_fdb_cap,omitempty"`
	LinearFDBTop         *int   `json:"linear_fdb_top,omitempty"`
	DefPort              *int   `json:"def_port,omitempty"`
	DefMCastPriPort      *int   `json:"def_mcast_pri_port,omitempty"`
	DefMCastNotPriPort   *int   `json:"def_mcast_not_pri_port,omitempty"`
	LifeTimeValue        *int   `json:"life_time_value,omitempty"`
	PortStateChange      *int   `json:"port_state_change,omitempty"`
	OptimizedSLVLMapping *int   `json:"optimized_slvl_mapping,omitempty"`
	LidsPerPort          *int   `json:"lids_per_port,omitempty"`
	PartEnfCap           *int   `json:"part_enf_cap,omitempty"`
	InbEnfCap            *int   `json:"inb_enf_cap,omitempty"`
	OutbEnfCap           *int   `json:"outb_enf_cap,omitempty"`
	FilterRawInbCap      *int   `json:"filter_raw_inb_cap,omitempty"`
	FilterRawOutbCap     *int   `json:"filter_raw_outb_cap,omitempty"`
	ENP0                 *int   `json:"enp0,omitempty"`
	MCastFDBTop          *int   `json:"mcast_fdb_top,omitempty"`
}
