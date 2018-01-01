package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type CoolerNode interface {
	graph.Node
	nodes.ComplexGasChannel
}

func NewCoolerNode(tOut float64, sigma float64) CoolerNode {
	var result = &coolerNode{
		tOut:  tOut,
		sigma: sigma,
	}

	result.complexGasInput = graph.NewAttachedPort(result)
	result.complexGasOutput = graph.NewAttachedPort(result)

	return result
}

type coolerNode struct {
	graph.BaseNode

	complexGasInput  graph.Port
	complexGasOutput graph.Port

	tOut  float64
	sigma float64
}

func (node *coolerNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Cooler")
}

func (node *coolerNode) GetPorts() []graph.Port {
	return []graph.Port{node.complexGasInput, node.complexGasOutput}
}

func (node *coolerNode) GetRequirePorts() []graph.Port {
	return []graph.Port{node.complexGasInput}
}

func (node *coolerNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{node.complexGasOutput}
}

func (node *coolerNode) Process() error {
	var state = node.complexGasInput.GetState().(states.ComplexGasPortState)
	state.TStag = node.tOut
	state.PStag *= node.sigma
	node.complexGasOutput.SetState(state)
	return nil
}

func (node *coolerNode) ComplexGasInput() graph.Port {
	return node.complexGasInput
}

func (node *coolerNode) ComplexGasOutput() graph.Port {
	return node.complexGasOutput
}
