package domain

type TopologyGroup struct {
	Key       string   `json:"key"`
	NodeGUIDs []string `json:"node_guids"`
}

type TopologyResponse struct {
	Nodes  []Node          `json:"nodes"`
	Groups []TopologyGroup `json:"groups"`
}
