package core

func graphErrorFromNodes(msg string, nodes []Node) GraphError {
	return &graphErrorImpl{
		msg:   msg,
		nodes: nodes,
	}
}

func graphErrorFromPorts(msg string, ports []Port) GraphError {
	return &graphErrorImpl{
		msg:   msg,
		ports: ports,
	}
}

type GraphError interface {
	error
	Nodes() []Node
	Ports() []Port
}

type graphErrorImpl struct {
	msg   string
	nodes []Node
	ports []Port
}

func (e *graphErrorImpl) Error() string {
	return e.msg
}

func (e *graphErrorImpl) Nodes() []Node {
	return e.nodes
}

func (e *graphErrorImpl) Ports() []Port {
	return e.ports
}
