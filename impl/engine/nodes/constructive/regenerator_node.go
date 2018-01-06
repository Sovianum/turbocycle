package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
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
	result.baseRegenerator = newBaseRegenerator(result)

	return result
}

type regeneratorNode struct {
	*baseRegenerator

	sigma     float64
	precision float64
}

func (node *regeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Regenerator")
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
