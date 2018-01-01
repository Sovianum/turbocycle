package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
)

type GasSinkNode interface {
	graph.Node
	nodes.GasSink
}

func NewGasSinkNode() GasSinkNode {
	var result = &gasSinkNode{}
	result.setInput(graph.NewAttachedPort(result))
	return result
}

type gasSinkNode struct {
	sinkNode
}

func (node *gasSinkNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "GasSink")
}

func (node *gasSinkNode) GasInput() graph.Port {
	return node.getInput()
}
