package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
)

type PowerSinkNode interface {
	graph.Node
	nodes.PowerSink
}

func NewPowerSinkNode() PowerSinkNode {
	var result = &powerSinkNode{}
	result.setInput(graph.NewAttachedPort(result))
	return result
}

type powerSinkNode struct {
	sinkNode
}

func (node *powerSinkNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "PowerSink")
}

func (node *powerSinkNode) PowerInput() graph.Port {
	return node.getInput()
}
