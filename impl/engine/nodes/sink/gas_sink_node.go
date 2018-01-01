package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
)

type GasSinkNode interface {
	core.Node
	nodes.GasSink
}

func NewGasSinkNode() GasSinkNode {
	var result = &gasSinkNode{}
	result.setInput(core.NewAttachedPort(result))
	return result
}

type gasSinkNode struct {
	sinkNode
}

func (node *gasSinkNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "GasSink")
}

func (node *gasSinkNode) GasInput() core.Port {
	return node.getInput()
}
