package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type MassRateRelSinkNode interface {
	core.Node
	nodes.MassRateRelSink
}

func NewMassRateRelSinkNode() MassRateRelSinkNode {
	var result = &massRateRelSinkNode{}
	result.setInput(core.NewAttachedPort(result))
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

func (node *massRateRelSinkNode) MassRateRelInput() core.Port {
	return node.getInput()
}
