package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type CoolerNode interface {
	core.Node
	nodes.ComplexGasChannel
}

func NewCoolerNode(tOut float64, sigma float64) CoolerNode {
	var result = &coolerNode{
		tOut:  tOut,
		sigma: sigma,
	}

	result.complexGasInput = core.NewAttachedPort(result)
	result.complexGasOutput = core.NewAttachedPort(result)

	return result
}

type coolerNode struct {
	core.BaseNode

	complexGasInput  core.Port
	complexGasOutput core.Port

	tOut  float64
	sigma float64
}

func (node *coolerNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Cooler")
}

func (node *coolerNode) GetPorts() []core.Port {
	return []core.Port{node.complexGasInput, node.complexGasOutput}
}

func (node *coolerNode) GetRequirePorts() []core.Port {
	return []core.Port{node.complexGasInput}
}

func (node *coolerNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.complexGasOutput}
}

func (node *coolerNode) Process() error {
	var state = node.complexGasInput.GetState().(states.ComplexGasPortState)
	state.TStag = node.tOut
	state.PStag *= node.sigma
	node.complexGasOutput.SetState(state)
	return nil
}

func (node *coolerNode) ComplexGasInput() core.Port {
	return node.complexGasInput
}

func (node *coolerNode) ComplexGasOutput() core.Port {
	return node.complexGasOutput
}
