package constructive

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

func newBaseCompressor(compressorNode graph.Node, precision float64) *baseCompressor {
	var result = &baseCompressor{
		precision: precision,
	}

	// ports are attached to compressorNode instead of result cos baseCompressor is not a node
	graph.AttachAllWithTags(
		compressorNode,
		[]*graph.Port{
			&result.powerOutput,
			&result.gasInput, &result.temperatureInput, &result.pressureInput,
			&result.gasOutput, &result.temperatureOutput, &result.pressureOutput, &result.massRateOutput,
		},
		[]string{
			nodes.PowerOutputTag,
			nodes.GasInputTag, nodes.TemperatureInputTag, nodes.PressureInputTag,
			nodes.GasOutputTag, nodes.TemperatureOutputTag, nodes.PowerOutputTag, nodes.MassRateOutputTag,
		},
	)

	return result
}

// TODO add collector port
type baseCompressor struct {
	graph.BaseNode

	gasInput         graph.Port
	temperatureInput graph.Port
	pressureInput    graph.Port

	gasOutput         graph.Port
	temperatureOutput graph.Port
	pressureOutput    graph.Port
	massRateOutput    graph.Port

	powerOutput graph.Port

	precision float64
}

func (node *baseCompressor) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *baseCompressor) GasInput() graph.Port {
	return node.gasInput
}

func (node *baseCompressor) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *baseCompressor) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *baseCompressor) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *baseCompressor) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *baseCompressor) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *baseCompressor) GetPorts() []graph.Port {
	return []graph.Port{
		node.powerOutput,
		node.gasInput, node.temperatureInput, node.pressureInput,
		node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
	}
}

func (node *baseCompressor) GetRequirePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.gasInput, node.temperatureInput, node.pressureInput,
	}, nil
}

func (node *baseCompressor) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.powerOutput,
		node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
	}, nil
}

func (node *baseCompressor) PowerOutput() graph.Port {
	return node.powerOutput
}

func (node *baseCompressor) TStagIn() float64 {
	return node.tStagIn()
}

func (node *baseCompressor) TStagOut() float64 {
	return node.tStagOut()
}

func (node *baseCompressor) PStagIn() float64 {
	return node.pStagIn()
}

func (node *baseCompressor) PStagOut() float64 {
	return node.pStagOut()
}

func (node *baseCompressor) LSpecific() float64 {
	return node.lSpecific()
}

func (node *baseCompressor) lSpecific() float64 {
	var cpMean = gases.CpMean(node.gas(), node.tStagIn(), node.tStagOut(), nodes.DefaultN)
	return cpMean * (node.tStagOut() - node.tStagIn())
}

func (node *baseCompressor) getTStagOut(tStagOutInit, piStag, etaAd float64) (float64, error) {
	var k = gases.K(node.gas(), node.tStagIn())
	iterFunc := func(t float64) (float64, error) {
		//k = gases.KMean(node.gas(), node.tStagIn(), t, turbine.DefaultN)
		x := math.Pow(piStag, (k-1)/k)
		return node.tStagIn() * (1 + (x-1)/etaAd), nil
	}
	k0 := 1.4
	x0 := math.Pow(piStag, (k0-1)/k0)

	t, e := common.SolveIteratively(iterFunc, node.tStagIn()*(1+(x0-1)/etaAd), node.precision, 1, 1000)
	return t, e
}

func (node *baseCompressor) tStagOutNewFunc(tStagOutCurr, piStag, etaAd float64) float64 {
	var x = node.xFunc(tStagOutCurr, piStag)
	return node.tStagIn() * (1 + (x-1)/etaAd)
}

func (node *baseCompressor) xFunc(tStagOut, piStag float64) float64 {
	var k = gases.KMean(node.gas(), node.tStagIn(), tStagOut, nodes.DefaultN)
	return math.Pow(piStag, (k-1)/k)
}

func (node *baseCompressor) tStagIn() float64 {
	return node.temperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseCompressor) tStagOut() float64 {
	return node.temperatureOutput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseCompressor) pStagIn() float64 {
	return node.pressureInput.GetState().(states.PressurePortState).PStag
}

func (node *baseCompressor) pStagOut() float64 {
	return node.pressureOutput.GetState().(states.PressurePortState).PStag
}

func (node *baseCompressor) gas() gases.Gas {
	return node.gasInput.GetState().(states.GasPortState).Gas
}
