package domain

import "fmt"

type NodeType int

const (
	NodeTypeCA     NodeType = 1 // host / channel adapter
	NodeTypeSwitch NodeType = 2 // switch
	NodeTypeRouter NodeType = 3 // router
)

func (t NodeType) String() string {
	switch t {
	case NodeTypeCA:
		return "host"
	case NodeTypeSwitch:
		return "switch"
	case NodeTypeRouter:
		return "router"
	default:
		return fmt.Sprintf("unknown(%d)", int(t))
	}
}

type Node struct {
	ID              int      `json:"id"`
	LogID           int      `json:"log_id"`
	NodeDesc        string   `json:"node_desc"`
	NumPorts        int      `json:"num_ports"`
	NodeType        NodeType `json:"node_type"`
	SystemImageGUID string   `json:"system_image_guid"`
	NodeGUID        string   `json:"node_guid"`
	PortGUID        string   `json:"port_guid"`
}
