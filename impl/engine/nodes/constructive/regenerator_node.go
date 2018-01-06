package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type RegeneratorNode interface {
	graph.Node

	HotInput() nodes.ComplexGasSink
	HotOutput() nodes.ComplexGasSource

	ColdInput() nodes.ComplexGasSink
	ColdOutput() nodes.ComplexGasSource

	Sigma() float64
}

func NewRegeneratorNode(sigma, precision float64) RegeneratorNode {
	var result = &regeneratorNode{
		sigma:     sigma,
		precision: precision,
	}

	graph.AttachAllPorts(
		result,
		&result.hotTemperatureInput, &result.hotPressureInput, &result.hotGasInput, &result.hotMassRateInput,
		&result.hotTemperatureOutput, &result.hotPressureOutput, &result.hotGasOutput, &result.hotMassRateOutput,
		&result.coldTemperatureInput, &result.coldPressureInput, &result.coldGasInput, &result.coldMassRateInput,
		&result.coldTemperatureOutput, &result.coldPressureOutput, &result.coldGasOutput, &result.coldMassRateOutput,
	)

	result.hotGasInput = graph.NewAttachedPort(result)
	result.coldGasInput = graph.NewAttachedPort(result)
	result.hotGasOutput = graph.NewAttachedPort(result)
	result.coldGasOutput = graph.NewAttachedPort(result)

	return result
}

type regeneratorNode struct {
	graph.BaseNode

	hotTemperatureInput graph.Port
	hotPressureInput    graph.Port
	hotGasInput         graph.Port
	hotMassRateInput    graph.Port

	hotTemperatureOutput graph.Port
	hotPressureOutput    graph.Port
	hotGasOutput         graph.Port
	hotMassRateOutput    graph.Port

	coldTemperatureInput graph.Port
	coldPressureInput    graph.Port
	coldGasInput         graph.Port
	coldMassRateInput    graph.Port

	coldTemperatureOutput graph.Port
	coldPressureOutput    graph.Port
	coldGasOutput         graph.Port
	coldMassRateOutput    graph.Port

	sigma     float64
	precision float64
}

func (node *regeneratorNode) HotInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.hotGasInput, node.hotTemperatureInput, node.hotPressureInput, node.hotMassRateInput,
	)
}

func (node *regeneratorNode) HotOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.hotGasOutput, node.hotTemperatureOutput, node.hotPressureOutput, node.hotMassRateOutput,
	)
}

func (node *regeneratorNode) ColdInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.coldGasInput, node.coldTemperatureInput, node.coldPressureInput, node.coldMassRateInput,
	)
}

func (node *regeneratorNode) ColdOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.coldGasOutput, node.coldTemperatureOutput, node.coldPressureOutput, node.coldMassRateOutput,
	)
}

func (node *regeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Regenerator")
}

func (node *regeneratorNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.hotTemperatureInput, node.hotPressureInput, node.hotGasInput, node.hotMassRateInput,
		node.hotTemperatureOutput, node.hotPressureOutput, node.hotGasOutput, node.hotMassRateOutput,
		node.coldTemperatureInput, node.coldPressureInput, node.coldGasInput, node.coldMassRateInput,
		node.coldTemperatureOutput, node.coldPressureOutput, node.coldGasOutput, node.coldMassRateOutput,
	}
}

func (node *regeneratorNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.hotTemperatureInput, node.hotPressureInput, node.hotGasInput, node.hotMassRateInput,
		node.coldTemperatureInput, node.coldPressureInput, node.coldGasInput, node.coldMassRateInput,
	}
}

func (node *regeneratorNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.hotTemperatureOutput, node.hotPressureOutput, node.hotGasOutput, node.hotMassRateOutput,
		node.coldTemperatureOutput, node.coldPressureOutput, node.coldGasOutput, node.coldMassRateOutput,
	}
}

func (node *regeneratorNode) Sigma() float64 {
	return node.sigma
}

func (node *regeneratorNode) Process() error {
	var iterFunc = node.getNewTOutIter

	var tStagColdOut, tStagHotOut = node.getNewTOut(node.tStagColdIn(), node.tStagHotIn(), iterFunc)
	graph.SetAll(
		[]graph.PortState{
			states.NewTemperaturePortState(tStagColdOut), states.NewTemperaturePortState(tStagHotOut),
			node.coldMassRateInput.GetState(), node.hotMassRateInput.GetState(),
			node.coldGasInput.GetState(), node.hotGasInput.GetState(),
			node.coldPressureInput.GetState(), node.hotPressureInput.GetState(),
		},
		[]graph.Port{
			node.coldTemperatureOutput, node.hotTemperatureOutput,
			node.coldMassRateOutput, node.hotMassRateOutput,
			node.coldGasOutput, node.hotGasOutput,
			node.coldPressureOutput, node.hotPressureOutput,
		},
	)

	graph.SetAll(
		[]graph.PortState{
			states.NewTemperaturePortState(tStagColdOut), states.NewTemperaturePortState(tStagHotOut),
		},
		[]graph.Port{
			node.coldTemperatureOutput, node.hotTemperatureOutput,
		},
	)

	return nil
}

func (node *regeneratorNode) getNewTOut(
	tStagColdOutCurr, tStagHotOutCurr float64,
	iterFunc func(float64, float64) (float64, float64),
) (tStagColdOutNew, tStagHotOutNew float64) {
	tStagColdOutNew, tStagHotOutNew = iterFunc(tStagColdOutCurr, tStagHotOutCurr)

	for !(common.Converged(tStagColdOutCurr, tStagColdOutNew, node.precision) &&
		common.Converged(tStagHotOutCurr, tStagHotOutNew, node.precision)) {
		tStagColdOutCurr, tStagHotOutCurr = tStagColdOutNew, tStagHotOutNew
		tStagColdOutNew, tStagHotOutNew = iterFunc(tStagColdOutCurr, tStagHotOutCurr)
	}
	return
}

func (node *regeneratorNode) getNewTOutIter(tStagColdOutCurr, tStagHotOutCurr float64) (tStagColdOut, tStagHotOut float64) {
	var hotMassRate = node.hotMassRateInput.GetState().(states.MassRatePortState).MassRate
	var hotGas = node.hotGasInput.GetState().(states.GasPortState).Gas

	var coldMassRate = node.coldMassRateInput.GetState().(states.MassRatePortState).MassRate
	var coldGas = node.coldGasInput.GetState().(states.GasPortState).Gas

	var hotHeatRate = hotMassRate * gases.CpMean(hotGas, node.tStagHotIn(), tStagHotOutCurr, nodes.DefaultN)
	var coldHeatRate = coldMassRate * gases.CpMean(coldGas, node.tStagColdIn(), tStagColdOutCurr, nodes.DefaultN)
	var heatRateFactor = hotHeatRate / coldHeatRate

	tStagColdOut = node.tStagColdIn() + node.sigma*(node.tStagHotIn()-node.tStagColdIn())
	tStagHotOut = node.sigma/heatRateFactor*node.tStagColdIn() + (1-node.sigma/heatRateFactor)*node.tStagHotIn()
	return
}

func (node *regeneratorNode) tStagHotIn() float64 {
	return node.hotTemperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *regeneratorNode) tStagHotOut() float64 {
	return node.hotTemperatureOutput.GetState().(states.TemperaturePortState).TStag
}

func (node *regeneratorNode) tStagColdIn() float64 {
	return node.coldTemperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *regeneratorNode) tStagColdOut() float64 {
	return node.coldTemperatureOutput.GetState().(states.TemperaturePortState).TStag
}
