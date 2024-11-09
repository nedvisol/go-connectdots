package graphdb

type NodeInfo struct {
	Label string
	Id    string
	Attrs *map[string]interface{}
}

type EdgeInfo struct {
	Label string
	Id    string
	Attrs *map[string]interface{}
	Left  *NodeInfo
	Right *NodeInfo
}

type GraphDbService interface {
	CreateNode(node *NodeInfo) error
	UpdateNode(node *NodeInfo, allowUpsert bool) error
	DeleteNode(node *NodeInfo) error
	UpdateEdge(edge *EdgeInfo, allowUpsert bool) error
}
