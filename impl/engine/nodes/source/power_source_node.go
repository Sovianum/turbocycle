package source

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type PowerSourceNode interface {
	graph.Node
	nodes.PowerSource
}

func NewPowerSourceNode(lRel float64) PowerSourceNode {
	var result = &powerSourceNode{
		lSpecific: lRel,
	}

	result.setOutput(graph.NewAttachedPort(result))
	return result
}

type powerSourceNode struct {
	sourceNode
	lSpecific float64
}

func (node *powerSourceNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "PowerSource")
}

func (node *powerSourceNode) Process() error {
	node.getOutput().SetState(states.NewPowerPortState(node.lSpecific))
	return nil
}

func (node *powerSourceNode) PowerOutput() graph.Port {
	return node.getOutput()
}
