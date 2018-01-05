package constructive

import (
	"errors"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

func NewBaseBlockedTurbine(node graph.Node, precision float64) *baseBlockedTurbine {
	var result = &baseBlockedTurbine{
		precision: precision,
	}
	graph.AttachAllPorts(
		node,
		&result.powerInput, &result.powerOutput,
		&result.temperatureInput, &result.pressureInput, &result.gasInput, &result.massRateInput,
		&result.temperatureOutput, &result.pressureOutput, &result.gasOutput, &result.massRateOutput,
	)
	return result
}

type baseBlockedTurbine struct {
	graph.BaseNode

	powerInput  graph.Port
	powerOutput graph.Port

	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port
	massRateInput    graph.Port

	temperatureOutput graph.Port
	pressureOutput    graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port

	precision float64
}

func (node *baseBlockedTurbine) GetPorts() []graph.Port {
	return []graph.Port{
		node.powerInput, node.powerOutput,
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *baseBlockedTurbine) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.powerInput, node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
	}
}

func (node *baseBlockedTurbine) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.powerOutput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *baseBlockedTurbine) LSpecific() float64 {
	return node.turbineLabour()
}

func (node *baseBlockedTurbine) MassRateRel() float64 {
	return node.massRateOutput.GetState().(states.MassRatePortState).MassRate
}

func (node *baseBlockedTurbine) InputGas() gases.Gas {
	return node.inputGas()
}

func (node *baseBlockedTurbine) TStagIn() float64 {
	return node.tStagIn()
}

func (node *baseBlockedTurbine) PStagIn() float64 {
	return node.pStagIn()
}

func (node *baseBlockedTurbine) TStagOut() float64 {
	return node.tStagOut()
}

func (node *baseBlockedTurbine) PStagOut() float64 {
	return node.pStagOut()
}

func (node *baseBlockedTurbine) PowerInput() graph.Port {
	return node.powerInput
}

func (node *baseBlockedTurbine) PowerOutput() graph.Port {
	return node.powerOutput
}

func (node *baseBlockedTurbine) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *baseBlockedTurbine) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *baseBlockedTurbine) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *baseBlockedTurbine) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *baseBlockedTurbine) GasInput() graph.Port {
	return node.gasInput
}

func (node *baseBlockedTurbine) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *baseBlockedTurbine) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *baseBlockedTurbine) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *baseBlockedTurbine) getTStagOut(etaT float64) (float64, error) {
	var tTStagCurr = node.getInitTtStag(etaT)
	var tTStagNew = node.getNewTtStag(tTStagCurr, etaT)

	for !common.Converged(tTStagCurr, tTStagNew, node.precision) {
		if math.IsNaN(tTStagCurr) || math.IsNaN(tTStagNew) {
			return 0, errors.New("failed to converge: try different initial guess")
		}
		tTStagCurr = tTStagNew
		tTStagNew = node.getNewTtStag(tTStagCurr, etaT)
	}

	return tTStagNew, nil
}

func (node *baseBlockedTurbine) getInitTtStag(etaT float64) float64 {
	return node.getNewTtStag(0.8*node.tStagIn(), etaT) // TODO move 0.8 out of code
}

func (node *baseBlockedTurbine) getNewTtStag(currTtStag, etaT float64) float64 {
	var k = gases.KMean(node.inputGas(), node.tStagIn(), currTtStag, nodes.DefaultN)
	var piTStag = node.piTStag(currTtStag, etaT)

	return node.tStagIn() * (1 - (1-math.Pow(piTStag, (1-k)/k))*etaT)
}

func (node *baseBlockedTurbine) turbineLabour() float64 {
	return -node.powerInput.GetState().(states.PowerPortState).LSpecific
}

func (node *baseBlockedTurbine) tStatOut(lambdaOut float64) float64 {
	var tStagOut = node.tStagOut()
	var k = gases.K(node.inputGas(), tStagOut)
	return tStagOut * gdf.Tau(lambdaOut, k)
}

func (node *baseBlockedTurbine) pStatOut(lambdaOut float64) float64 {
	var pStagOut = node.pStagOut()
	var tStagOut = node.tStagOut()
	var k = gases.K(node.inputGas(), tStagOut)
	return pStagOut * gdf.Tau(lambdaOut, k)
}

func (node *baseBlockedTurbine) inputGas() gases.Gas {
	return node.gasInput.GetState().(states.GasPortState).Gas
}

func (node *baseBlockedTurbine) tStagIn() float64 {
	return node.temperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseBlockedTurbine) pStagIn() float64 {
	return node.pressureInput.GetState().(states.PressurePortState).PStag
}

func (node *baseBlockedTurbine) tStagOut() float64 {
	return node.temperatureOutput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseBlockedTurbine) pStagOut() float64 {
	return node.pressureOutput.GetState().(states.PressurePortState).PStag
}
