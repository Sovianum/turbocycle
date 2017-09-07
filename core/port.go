package core

type Port struct {
	state     IPortState
	innerNode Node
	outerNode Node
	linkPort  *Port
}

func NewPort() *Port {
	return &Port{
		state:     nil,
		innerNode: nil,
		outerNode: nil,
		linkPort:  nil,
	}
}

func (port *Port) GetState() IPortState {
	return port.state
}

func (port *Port) SetState(state IPortState) {
	port.state = state
}

func (port *Port) GetInnerNode() Node {
	return port.innerNode
}

func (port *Port) SetInnerNode(src Node) {
	port.innerNode = src
}

func (port *Port) GetOuterNode() Node {
	return port.outerNode
}

func (port *Port) SetOuterNode(dest Node) {
	port.outerNode = dest
}

func Link(port1 *Port, port2 *Port) {
	port1.outerNode = port2.innerNode
	port2.outerNode = port1.innerNode
}
