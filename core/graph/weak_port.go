package graph

func NewTransformWeakPort(referencePort Port, transformFunc func(PortState) PortState) Port {
	return &weakPort{
		referencePort: referencePort,
		transformFunc: transformFunc,
	}
}

func NewWeakPort(referencePort Port) Port {
	return &weakPort{referencePort: referencePort}
}

type weakPort struct {
	referencePort Port
	transformFunc func(PortState) PortState

	outerNode Node
	linkPort  Port
}

func (p *weakPort) GetTag() string {
	return p.referencePort.GetTag()
}

func (p *weakPort) SetTag(tag string) {
	p.referencePort.SetTag(tag)
}

func (p *weakPort) GetState() PortState {
	if p.transformFunc == nil {
		return p.referencePort.GetState()
	}
	return p.transformFunc(p.referencePort.GetState())
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
