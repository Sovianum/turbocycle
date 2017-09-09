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

type port struct {
	state     IPortState
	innerNode Node
	outerNode Node
	linkPort  Port
}

func NewPort() Port {
	return &port{
		state:     nil,
		innerNode: nil,
		outerNode: nil,
		linkPort:  nil,
	}
}

func Link(port1 Port, port2 Port) {
	port1.SetOuterNode(port2.GetInnerNode())
	port2.SetOuterNode(port1.GetInnerNode())
}

func (port *port) GetState() IPortState {
	return port.state
}

func (port *port) SetState(state IPortState) {
	port.state = state
}

func (port *port) GetInnerNode() Node {
	return port.innerNode
}

func (port *port) SetInnerNode(src Node) {
	port.innerNode = src
}

func (port *port) GetOuterNode() Node {
	return port.outerNode
}

func (port *port) SetOuterNode(dest Node) {
	port.outerNode = dest
}

func (port *port) GetLinkPort() Port {
	return port.linkPort
}

func (port *port) SetLinkPort(another Port) {
	port.linkPort = another
}
