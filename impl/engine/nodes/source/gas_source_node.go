package source

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type GasSourceNode interface {
	core.Node
	nodes.GasSource
}

func NewGasSourceNode(gas gases.Gas) GasSourceNode {
	var result = &gasSourceNode{
		gas: gas,
	}
	result.setOutput(core.NewAttachedPort(result))
	return result
}

type gasSourceNode struct {
	sourceNode
	gas gases.Gas
}

func (node *gasSourceNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "GasSource")
}

func (node *gasSourceNode) Process() error {
	node.output.SetState(states.NewGasPortState(node.gas))
	return nil
}

func (node *gasSourceNode) GasOutput() core.Port {
	return node.getOutput()
}
