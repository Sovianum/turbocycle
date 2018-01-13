package source

import "github.com/Sovianum/turbocycle/core/graph"

type sourceNode struct {
	graph.BaseNode
	output graph.Port
}

func (node *sourceNode) GetRequirePorts() ([]graph.Port, error) {
	return make([]graph.Port, 0), nil
}

func (node *sourceNode) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{node.output}, nil
}

func (node *sourceNode) GetPorts() []graph.Port {
	return []graph.Port{node.output}
}

func (node *sourceNode) setOutput(output graph.Port) {
	node.output = output
}

func (node *sourceNode) getOutput() graph.Port {
	return node.output
}
