package source

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type MassRateRelSourceNode interface {
	core.Node
	nodes.MassRateRelSource
}

func NewMassRateRelSourceNode(massRateRel float64) MassRateRelSourceNode {
	var result = &massRateRelSourceNode{
		massRateRel: massRateRel,
	}

	result.setOutput(core.NewAttachedPort(result))
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
	node.getOutput().SetState(states.NewMassRateRelPortState(node.massRateRel))
	return nil
}

func (node *massRateRelSourceNode) MassRateRelOut() float64 {
	return node.massRateRel
}

func (node *massRateRelSourceNode) MassRateRelOutput() core.Port {
	return node.getOutput()
}
