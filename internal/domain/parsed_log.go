package domain

type ParsedLog struct {
	Log       Log
	Nodes     []Node
	Ports     []Port
	NodesInfo []NodeInfo
	Topology  []TopologyGroup
}
