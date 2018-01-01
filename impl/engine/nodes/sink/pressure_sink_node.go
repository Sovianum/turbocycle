package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type PressureSinkNode interface {
	core.Node
	nodes.PressureSink
}

func NewPressureSinkNode() PressureSinkNode {
	var result = &pressureSinkNode{}
	result.setInput(core.NewAttachedPort(result))
	return result
}

type pressureSinkNode struct {
	sinkNode
}

func (node *pressureSinkNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "PressureSink")
}

func (node *pressureSinkNode) PStagIn() float64 {
	return node.getInput().GetState().(states.PressurePortState).PStag
}

func (node *pressureSinkNode) PressureInput() core.Port {
	return node.getInput()
}
