package constructive

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
)

func newBaseBurner(node graph.Node, fuel fuel.GasFuel, etaBurn, tFuel, t0, precision float64) *baseBurner {
	var result = &baseBurner{
		fuel:      fuel,
		etaBurn:   etaBurn,
		tFuel:     tFuel,
		t0:        t0,
		precision: precision,
	}

	graph.AttachAllPorts(
		node,
		&result.temperatureInput, &result.pressureInput, &result.gasInput, &result.massRateInput,
		&result.temperatureOutput, &result.pressureOutput, &result.gasOutput, &result.massRateOutput,
	)
	return result
}

type baseBurner struct {
	graph.BaseNode

	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port
	massRateInput    graph.Port

	temperatureOutput graph.Port
	pressureOutput    graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port

	fuel      fuel.GasFuel
	tFuel     float64
	etaBurn   float64
	t0        float64
	precision float64
}

func (node *baseBurner) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *baseBurner) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *baseBurner) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *baseBurner) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *baseBurner) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *baseBurner) GasInput() graph.Port {
	return node.gasInput
}

func (node *baseBurner) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *baseBurner) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *baseBurner) TStagIn() float64 {
	return node.tStagIn()
}

func (node *baseBurner) TStagOut() float64 {
	return node.tStagOut()
}

func (node *baseBurner) PStagIn() float64 {
	return node.pStagIn()
}

func (node *baseBurner) PStagOut() float64 {
	return node.pStagOut()
}

func (node *baseBurner) GetPorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *baseBurner) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
	}
}

func (node *baseBurner) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *baseBurner) Fuel() fuel.GasFuel {
	return node.fuel
}

func (node *baseBurner) Eta() float64 {
	return node.etaBurn
}

func (node *baseBurner) T0() float64 {
	return node.t0
}

func (node *baseBurner) TFuel() float64 {
	return node.tFuel
}

func (node *baseBurner) inletGas() gases.Gas {
	return node.gasInput.GetState().(states.GasPortState).Gas
}

func (node *baseBurner) tStagIn() float64 {
	return node.temperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseBurner) tStagOut() float64 {
	return node.temperatureOutput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseBurner) pStagIn() float64 {
	return node.pressureInput.GetState().(states.PressurePortState).PStag
}

func (node *baseBurner) pStagOut() float64 {
	return node.pressureOutput.GetState().(states.PressurePortState).PStag
}
