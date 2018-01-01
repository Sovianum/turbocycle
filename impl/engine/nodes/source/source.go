package source

import (
	"github.com/Sovianum/turbocycle/core"
)

type sourceNode struct {
	core.BaseNode
	output core.Port
}

func (node *sourceNode) GetRequirePorts() []core.Port {
	return make([]core.Port, 0)
}

func (node *sourceNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.output}
}

func (node *sourceNode) GetPorts() []core.Port {
	return []core.Port{node.output}
}

func (node *sourceNode) setOutput(output core.Port) {
	node.output = output
}

func (node *sourceNode) getOutput() core.Port {
	return node.output
}
