package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type MassRateRelSinkNode interface {
	graph.Node
	nodes.MassRateRelSink
}

func NewMassRateRelSinkNode() MassRateRelSinkNode {
	var result = &massRateRelSinkNode{}
	result.setInput(graph.NewAttachedPort(result))
	return result
}

type massRateRelSinkNode struct {
	sinkNode
}

func (node *massRateRelSinkNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "MassRateRelSink")
}

func (node *massRateRelSinkNode) MassRateRelIn() float64 {
	return node.getInput().GetState().(states.MassRateRelPortState).MassRateRel
}

func (node *massRateRelSinkNode) MassRateRelInput() graph.Port {
	return node.getInput()
}
