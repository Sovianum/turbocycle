package core

func graphErrorFromNodes(msg string, nodes []Node) *graphError {
	return &graphError{
		msg:   msg,
		nodes: nodes,
	}
}

func graphErrorFromPorts(msg string, ports []Port) *graphError {
	return &graphError{
		msg:   msg,
		ports: ports,
	}
}

type graphError struct {
	msg   string
	nodes []Node
	ports []Port
}

func (e *graphError) Error() string {
	return e.msg
}

func (e *graphError) Nodes() []Node {
	return e.nodes
}

func (e *graphError) Ports() []Port {
	return e.ports
}
