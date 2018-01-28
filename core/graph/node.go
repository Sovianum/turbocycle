package graph

type Node interface {
	SetName(name string)
	GetName() string

	Process() error
	GetRequirePorts() ([]Port, error)
	GetUpdatePorts() ([]Port, error)
	GetPorts() []Port
	ContextDefined(key int) bool
}

func NewBaseNode() *BaseNode {
	return &BaseNode{}
}

type BaseNode struct {
	name string
}

func (node *BaseNode) SetName(name string) {
	node.name = name
}

func (node *BaseNode) ContextDefined(key int) bool {
	return true
}

func (node *BaseNode) GetInstanceName() string {
	return node.name
}
