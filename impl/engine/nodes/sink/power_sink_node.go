package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
)

type PowerSinkNode interface {
	core.Node
	nodes.PowerSink
}

func NewPowerSinkNode() PowerSinkNode {
	var result = &powerSinkNode{}
	result.setInput(core.NewAttachedPort(result))
	return result
}

type powerSinkNode struct {
	sinkNode
}

func (node *powerSinkNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "PowerSink")
}

func (node *powerSinkNode) PowerInput() core.Port {
	return node.getInput()
}
