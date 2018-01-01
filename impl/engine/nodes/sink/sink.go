package sink

import "github.com/Sovianum/turbocycle/core"

type sinkNode struct {
	core.BaseNode
	input core.Port
}

func (node *sinkNode) GetRequirePorts() []core.Port {
	return []core.Port{node.input}
}

func (node *sinkNode) GetUpdatePorts() []core.Port {
	return make([]core.Port, 0)
}

func (node *sinkNode) GetPorts() []core.Port {
	return []core.Port{node.input}
}

func (node *sinkNode) Process() error {
	return nil
}

func (node *sinkNode) setInput(input core.Port) {
	node.input = input
}

func (node *sinkNode) getInput() core.Port {
	return node.input
}
