package core

type Port interface {
	GetState() IPortState
	SetState(state IPortState)
	GetInnerNode() Node
	SetInnerNode(Node)
	GetOuterNode() Node
	SetOuterNode(Node)
	GetLinkPort() Port
	SetLinkPort(Port)
}

type portType struct {
	state     IPortState
	innerNode Node
	outerNode Node
	linkPort  Port
}

func NewPort() Port {
	return &portType{
		state:     nil,
		innerNode: nil,
		outerNode: nil,
		linkPort:  nil,
	}
}

func Link(port1 Port, port2 Port) {
	port1.SetLinkPort(port2)
	port2.SetLinkPort(port1)

	port1.SetOuterNode(port2.GetInnerNode())
	port2.SetOuterNode(port1.GetInnerNode())
}

func (port *portType) GetState() IPortState {
	return port.state
}

func (port *portType) SetState(state IPortState) {
	port.state = state
	if port.linkPort != nil {
		port.linkPort.(*portType).state = state
	}
}

func (port *portType) GetInnerNode() Node {
	return port.innerNode
}

func (port *portType) SetInnerNode(src Node) {
	port.innerNode = src
}

func (port *portType) GetOuterNode() Node {
	return port.outerNode
}

func (port *portType) SetOuterNode(dest Node) {
	port.outerNode = dest
}

func (port *portType) GetLinkPort() Port {
	return port.linkPort
}

func (port *portType) SetLinkPort(another Port) {
	port.linkPort = another
}
