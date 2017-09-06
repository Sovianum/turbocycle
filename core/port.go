package core

type Port struct {
	state    IPortState
	src      Node
	dest     Node
	linkPort *Port
}

func NewPort() *Port {
	return &Port{
		state:    nil,
		src:      nil,
		dest:     nil,
		linkPort: nil,
	}
}

func (port *Port) GetState() IPortState {
	return port.state
}

func (port *Port) SetState(state IPortState) {
	port.state = state
}

func (port *Port) GetSrc() Node {
	return port.src
}

func (port *Port) SetSrc(src Node) {
	port.src = src
}

func (port *Port) GetDest() Node {
	return port.dest
}

func (port *Port) SetDest(dest Node) {
	port.dest = dest
}

func Link(outputPort *Port, inputPort *Port) {
	inputPort.src = outputPort.src
	outputPort.dest = inputPort.dest

	outputPort.linkPort = inputPort
	inputPort.linkPort = outputPort
}
