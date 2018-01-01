package core

type Port interface {
	GetState() PortState
	SetState(state PortState)
	GetInnerNode() Node
	SetInnerNode(Node)
	GetOuterNode() Node
	SetOuterNode(Node)
	GetLinkPort() Port
	SetLinkPort(Port)
}

type portType struct {
	state     PortState
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

func NewAttachedPort(node Node) Port {
	var port = &portType{
		state:     nil,
		innerNode: nil,
		outerNode: nil,
		linkPort:  nil,
	}
	port.SetInnerNode(node)
	return port
}

func Link(port1 Port, port2 Port) {
	port1.SetLinkPort(port2)
	port2.SetLinkPort(port1)

	port1.SetOuterNode(port2.GetInnerNode())
	port2.SetOuterNode(port1.GetInnerNode())
}

func (port *portType) GetState() PortState {
	return port.state
}

func (port *portType) SetState(state PortState) {
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
