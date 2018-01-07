package graph

func NewWeakPort(referencePort Port) Port {
	return &weakPort{referencePort: referencePort}
}

type weakPort struct {
	referencePort Port

	outerNode Node
	linkPort  Port
}

func (p *weakPort) GetState() PortState {
	return p.referencePort.GetState()
}

func (p *weakPort) SetState(state PortState) {
	p.referencePort.SetState(state)
}

func (p *weakPort) GetInnerNode() Node {
	return p.referencePort.GetInnerNode()
}

func (p *weakPort) SetInnerNode(node Node) {
	p.referencePort.SetInnerNode(node)
}

func (p *weakPort) GetOuterNode() Node {
	return p.outerNode
}

func (p *weakPort) SetOuterNode(node Node) {
	p.outerNode = node
}

func (p *weakPort) GetLinkPort() Port {
	return p.linkPort
}

func (p *weakPort) SetLinkPort(port Port) {
	p.linkPort = port
}
