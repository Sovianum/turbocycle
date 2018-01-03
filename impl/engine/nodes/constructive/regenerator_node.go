package constructive

import (
	"fmt"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

const (
	SigmaByHotSide  = "sigmaByHotSide"
	SigmaByColdSide = "sigmaByColdSide"
)

type RegeneratorNode interface {
	graph.Node

	HotGasInput() graph.Port
	HotTemperatureInput() graph.Port
	HotPressureInput() graph.Port
	HotMassRateInput() graph.Port

	ColdGasInput() graph.Port
	ColdTemperatureInput() graph.Port
	ColdPressureInput() graph.Port
	ColdMassRateInput() graph.Port

	HotGasOutput() graph.Port
	HotTemperatureOutput() graph.Port
	HotPressureOutput() graph.Port
	HotMassRateOutput() graph.Port

	ColdGasOutput() graph.Port
	ColdTemperatureOutput() graph.Port
	ColdPressureOutput() graph.Port
	ColdMassRateOutput() graph.Port

	Sigma() float64
}

func NewRegeneratorNode(sigma, precision float64, mode string) RegeneratorNode {
	var result = &regeneratorNode{
		sigma:     sigma,
		precision: precision,
		mode:      mode,
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
	mode      string
}

func (node *regeneratorNode) HotGasInput() graph.Port {
	return node.hotGasInput
}

func (node *regeneratorNode) HotTemperatureInput() graph.Port {
	return node.hotTemperatureInput
}

func (node *regeneratorNode) HotPressureInput() graph.Port {
	return node.hotPressureInput
}

func (node *regeneratorNode) HotMassRateInput() graph.Port {
	return node.hotMassRateInput
}

func (node *regeneratorNode) ColdGasInput() graph.Port {
	return node.coldGasInput
}

func (node *regeneratorNode) ColdTemperatureInput() graph.Port {
	return node.coldTemperatureInput
}

func (node *regeneratorNode) ColdPressureInput() graph.Port {
	return node.coldPressureInput
}

func (node *regeneratorNode) ColdMassRateInput() graph.Port {
	return node.coldMassRateInput
}

func (node *regeneratorNode) HotGasOutput() graph.Port {
	return node.hotGasOutput
}

func (node *regeneratorNode) HotTemperatureOutput() graph.Port {
	return node.hotTemperatureOutput
}

func (node *regeneratorNode) HotPressureOutput() graph.Port {
	return node.hotPressureOutput
}

func (node *regeneratorNode) HotMassRateOutput() graph.Port {
	return node.hotMassRateOutput
}

func (node *regeneratorNode) ColdGasOutput() graph.Port {
	return node.coldGasOutput
}

func (node *regeneratorNode) ColdTemperatureOutput() graph.Port {
	return node.coldTemperatureOutput
}

func (node *regeneratorNode) ColdPressureOutput() graph.Port {
	return node.coldPressureOutput
}

func (node *regeneratorNode) ColdMassRateOutput() graph.Port {
	return node.coldMassRateOutput
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
	var iterFunc func(float64, float64) (float64, float64)
	switch node.mode {
	case SigmaByColdSide:
		iterFunc = node.getNewTOutSigmaByColdSide
	case SigmaByHotSide:
		iterFunc = node.getNewTOutSigmaByHotSide
	default:
		return fmt.Errorf("invalid Regenerator node state: %s", node.mode)
	}

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

func (node *regeneratorNode) getNewTOutSigmaByColdSide(tStagColdOutCurr, tStagHotOutCurr float64) (tStagColdOut, tStagHotOut float64) {
	var hotMassRate = node.hotMassRateInput.GetState().(states.MassRateRelPortState).MassRateRel
	var hotGas = node.hotGasInput.GetState().(states.GasPortState).Gas

	var coldMassRate = node.coldMassRateInput.GetState().(states.MassRateRelPortState).MassRateRel
	var coldGas = node.coldGasInput.GetState().(states.GasPortState).Gas

	var hotHeatRate = hotMassRate * gases.CpMean(hotGas, node.tStagHotIn(), tStagHotOutCurr, nodes.DefaultN)
	var coldHeatRate = coldMassRate * gases.CpMean(coldGas, node.tStagColdIn(), tStagColdOutCurr, nodes.DefaultN)
	var heatRateFactor = hotHeatRate / coldHeatRate

	tStagColdOut = node.tStagColdIn() + node.sigma*(node.tStagHotIn()-node.tStagColdIn())
	tStagHotOut = node.sigma/heatRateFactor*node.tStagColdIn() + (1-node.sigma/heatRateFactor)*node.tStagHotIn()
	return
}

func (node *regeneratorNode) getNewTOutSigmaByHotSide(tStagColdOutCurr, tStagHotOutCurr float64) (tStagColdOut, tStagHotOut float64) {
	var hotMassRate = node.hotMassRateInput.GetState().(states.MassRateRelPortState).MassRateRel
	var hotGas = node.hotGasInput.GetState().(states.GasPortState).Gas

	var coldMassRate = node.coldMassRateInput.GetState().(states.MassRateRelPortState).MassRateRel
	var coldGas = node.coldGasInput.GetState().(states.GasPortState).Gas

	var hotHeatRate = hotMassRate * gases.CpMean(hotGas, node.tStagHotIn(), tStagHotOutCurr, nodes.DefaultN)
	var coldHeatRate = coldMassRate * gases.CpMean(coldGas, node.tStagColdIn(), tStagColdOutCurr, nodes.DefaultN)
	var heatRateFactor = hotHeatRate / coldHeatRate

	tStagHotOut = node.tStagColdIn() + node.sigma*(node.tStagHotIn()-node.tStagColdIn())
	tStagColdOut = node.tStagColdIn() + heatRateFactor*node.sigma*(node.tStagHotIn()-node.tStagColdIn())
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
