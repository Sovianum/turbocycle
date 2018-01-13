package source

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type ComplexGasSourceNode interface {
	graph.Node
	nodes.ComplexGasSource
}

func NewComplexGasSourceNode(gas gases.Gas, tStag, pStag, massRate float64) ComplexGasSourceNode {
	var result = &complexGasSourceNode{
		gas:      gas,
		tStag:    tStag,
		pStag:    pStag,
		massRate: massRate,
	}
	graph.AttachAllWithTags(
		result,
		[]*graph.Port{
			&result.gasOutput, &result.temperatureOutput, &result.pressureOutput, &result.massRateOutput,
		},
		[]string{
			nodes.GasOutputTag, nodes.TemperatureOutputTag, nodes.PressureOutputTag, nodes.MassRateOutputTag,
		},
	)
	return result
}

type complexGasSourceNode struct {
	sourceNode

	gasOutput         graph.Port
	temperatureOutput graph.Port
	pressureOutput    graph.Port
	massRateOutput    graph.Port

	gas      gases.Gas
	tStag    float64
	pStag    float64
	massRate float64
}

func (node *complexGasSourceNode) GetRequirePorts() ([]graph.Port, error) {
	return make([]graph.Port, 0), nil
}

func (node *complexGasSourceNode) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
	}, nil
}

func (node *complexGasSourceNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
	}
}

func (node *complexGasSourceNode) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *complexGasSourceNode) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *complexGasSourceNode) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *complexGasSourceNode) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *complexGasSourceNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ComplexGasSource")
}

func (node *complexGasSourceNode) Process() error {
	node.gasOutput.SetState(states.NewGasPortState(node.gas))
	node.temperatureOutput.SetState(states.NewTemperaturePortState(node.tStag))
	node.pressureOutput.SetState(states.NewPressurePortState(node.pStag))
	node.massRateOutput.SetState(states.NewMassRatePortState(node.massRate))
	return nil
}
