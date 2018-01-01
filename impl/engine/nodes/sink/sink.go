package sink

import "github.com/Sovianum/turbocycle/core/graph"

type sinkNode struct {
	graph.BaseNode
	input graph.Port
}

func (node *sinkNode) GetRequirePorts() []graph.Port {
	return []graph.Port{node.input}
}

func (node *sinkNode) GetUpdatePorts() []graph.Port {
	return make([]graph.Port, 0)
}

func (node *sinkNode) GetPorts() []graph.Port {
	return []graph.Port{node.input}
}

func (node *sinkNode) Process() error {
	return nil
}

func (node *sinkNode) setInput(input graph.Port) {
	node.input = input
}

func (node *sinkNode) getInput() graph.Port {
	return node.input
}
