package source

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type MassRateRelSourceNode interface {
	graph.Node
	nodes.MassRateSource
}

func NewMassRateRelSourceNode(massRateRel float64) MassRateRelSourceNode {
	var result = &massRateRelSourceNode{
		massRateRel: massRateRel,
	}

	result.setOutput(graph.NewAttachedPort(result))
	return result
}

type massRateRelSourceNode struct {
	sourceNode
	massRateRel float64
}

func (node *massRateRelSourceNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "MassRateRelSource")
}

func (node *massRateRelSourceNode) Process() error {
	node.getOutput().SetState(states.NewMassRatePortState(node.massRateRel))
	return nil
}

func (node *massRateRelSourceNode) MassRateOut() float64 {
	return node.massRateRel
}

func (node *massRateRelSourceNode) MassRateOutput() graph.Port {
	return node.getOutput()
}
