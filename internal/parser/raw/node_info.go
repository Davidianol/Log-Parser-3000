package raw

type SwitchInfo struct {
	NodeGUID             string
	LinearFDBCap         string
	RandomFDBCap         string
	MCastFDBCap          string
	LinearFDBTop         string
	DefPort              string
	DefMCastPriPort      string
	DefMCastNotPriPort   string
	LifeTimeValue        string
	PortStateChange      string
	OptimizedSLVLMapping string
	LidsPerPort          string
	PartEnfCap           string
	InbEnfCap            string
	OutbEnfCap           string
	FilterRawInbCap      string
	FilterRawOutbCap     string
	ENP0                 string
	MCastFDBTop          string
}

type GeneralInfo struct {
	NodeGUID     string
	SerialNumber string
	PartNumber   string
	Revision     string
	ProductName  string
}

type SharpInfo struct {
	NodeGUID               string
	Endianness             string
	EnableEndiannessPerJob string
	ReproducibilityDisable string
}
