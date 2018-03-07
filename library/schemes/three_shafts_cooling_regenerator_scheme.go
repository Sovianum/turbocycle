package schemes

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

func NewThreeShaftsCoolingRegeneratorScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	middlePressureCascade compose.TurboCascadeNode,
	cooler constructive.CoolerNode,
	regenerativeGasGenerator compose.RegenerativeGasGeneratorNode,
	middlePressureCompressorPipe constructive.PressureLossNode,
	highPressureTurbinePipe constructive.PressureLossNode,
	middlePressureTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
) ThreeShaftsCoolingRegeneratorScheme {
	return &threeShaftsCoolingRegeneratorScheme{
		gasSource:                    gasSource,
		inletPressureDrop:            inletPressureDrop,
		middlePressureCascade:        middlePressureCascade,
		regenerativeGasGenerator:     regenerativeGasGenerator,
		middlePressureCompressorPipe: middlePressureCompressorPipe,
		cooler:                    cooler,
		highPressureTurbinePipe:   highPressureTurbinePipe,
		middlePressureTurbinePipe: middlePressureTurbinePipe,
		freeTurbineBlock:          freeTurbineBlock,

		breaker1: helper.NewComplexCycleBreakNode(
			gases.GetAir(), states.StandardTemperature, states.StandardPressure, 1,
		),
		breaker2: helper.NewComplexCycleBreakNode(
			gases.GetAir(), states.StandardTemperature, states.StandardPressure, 1,
		),
	}
}

type ThreeShaftsCoolingRegeneratorScheme interface {
	Scheme
	DoubleCompressor
	InitGasGeneratorCompressor(gas gases.Gas, tStag, pStag, massRate float64)
	InitGasGeneratorHeatExchanger(gas gases.Gas, tStag, pStag, massRate float64)
}

type threeShaftsCoolingRegeneratorScheme struct {
	gasSource                    source.ComplexGasSourceNode
	inletPressureDrop            constructive.PressureLossNode
	middlePressureCascade        compose.TurboCascadeNode
	cooler                       constructive.CoolerNode
	regenerativeGasGenerator     compose.RegenerativeGasGeneratorNode
	middlePressureCompressorPipe constructive.PressureLossNode
	highPressureTurbinePipe      constructive.PressureLossNode
	middlePressureTurbinePipe    constructive.PressureLossNode
	freeTurbineBlock             compose.FreeTurbineBlockNode

	gasSink         sink.SinkNode
	massRateSink    sink.SinkNode
	pressureSink    sink.SinkNode
	temperatureSink sink.SinkNode
	powerSink       sink.SinkNode

	breaker1 helper.ComplexCycleBreakNode
	breaker2 helper.ComplexCycleBreakNode
}

func (scheme *threeShaftsCoolingRegeneratorScheme) LPC() constructive.CompressorNode {
	return scheme.middlePressureCascade.Compressor()
}

func (scheme *threeShaftsCoolingRegeneratorScheme) HPC() constructive.CompressorNode {
	return scheme.regenerativeGasGenerator.TurboCascade().Compressor()
}

func (scheme *threeShaftsCoolingRegeneratorScheme) InitGasGeneratorCompressor(gas gases.Gas, tStag, pStag, massRate float64) {
	var b = scheme.breaker1
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gas), states.NewTemperaturePortState(tStag),
			states.NewPressurePortState(pStag), states.NewMassRatePortState(massRate),
		},
		[]graph.Port{
			b.GasInput(), b.TemperatureInput(), b.PressureInput(), b.MassRateInput(),
		},
	)
}

func (scheme *threeShaftsCoolingRegeneratorScheme) InitGasGeneratorHeatExchanger(gas gases.Gas, tStag, pStag, massRate float64) {
	var b = scheme.breaker2
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gas), states.NewTemperaturePortState(tStag),
			states.NewPressurePortState(pStag), states.NewMassRatePortState(massRate),
		},
		[]graph.Port{
			b.GasInput(), b.TemperatureInput(), b.PressureInput(), b.MassRateInput(),
		},
	)
}

func (scheme *threeShaftsCoolingRegeneratorScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return lSpecific * massRateRel
}

func (scheme *threeShaftsCoolingRegeneratorScheme) GetFuelMassRateRel() float64 {
	var burner = scheme.regenerativeGasGenerator.Burner()
	var massRateRel = burner.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return burner.FuelRateRel() * massRateRel
}

func (scheme *threeShaftsCoolingRegeneratorScheme) GetQLower() float64 {
	return scheme.regenerativeGasGenerator.Burner().Fuel().QLower()
}

func (scheme *threeShaftsCoolingRegeneratorScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()

	return graph.NewNetwork([]graph.Node{
		scheme.gasSource, scheme.inletPressureDrop, scheme.middlePressureCascade, scheme.cooler,
		scheme.middlePressureCompressorPipe, scheme.regenerativeGasGenerator, scheme.highPressureTurbinePipe,
		scheme.middlePressureTurbinePipe, scheme.freeTurbineBlock,
		scheme.gasSink, scheme.massRateSink, scheme.pressureSink, scheme.temperatureSink, scheme.powerSink,
		scheme.breaker1, scheme.breaker2,
	})
}

func (scheme *threeShaftsCoolingRegeneratorScheme) linkPorts() {
	nodes.LinkComplexOutToIn(scheme.gasSource, scheme.inletPressureDrop)
	nodes.LinkComplexOutToIn(scheme.inletPressureDrop, scheme.middlePressureCascade.CompressorComplexGasInput())

	nodes.LinkComplexOutToIn(scheme.middlePressureCascade.CompressorComplexGasOutput(), scheme.middlePressureCompressorPipe)
	nodes.LinkComplexOutToIn(scheme.middlePressureCompressorPipe, scheme.cooler)
	nodes.LinkComplexOutToIn(scheme.cooler, scheme.breaker1)
	nodes.LinkComplexOutToIn(scheme.breaker1, scheme.regenerativeGasGenerator.CompressorInput())
	nodes.LinkComplexOutToIn(scheme.regenerativeGasGenerator.HeatExchangerHotOutput(), scheme.highPressureTurbinePipe)
	nodes.LinkComplexOutToIn(scheme.highPressureTurbinePipe, scheme.middlePressureCascade.TurbineComplexGasInput())
	nodes.LinkComplexOutToIn(scheme.middlePressureCascade.TurbineComplexGasOutput(), scheme.middlePressureTurbinePipe)
	nodes.LinkComplexOutToIn(scheme.middlePressureTurbinePipe, scheme.freeTurbineBlock)

	nodes.LinkComplexOutToIn(scheme.freeTurbineBlock, scheme.breaker2)
	nodes.LinkComplexOutToIn(scheme.breaker2, scheme.regenerativeGasGenerator.HeatExchangerHotInput())

	scheme.gasSink = sink.SinkPort(scheme.regenerativeGasGenerator.TurbineOutput().GasOutput())
	scheme.temperatureSink = sink.SinkPort(scheme.regenerativeGasGenerator.TurbineOutput().TemperatureOutput())
	scheme.pressureSink = sink.SinkPort(scheme.regenerativeGasGenerator.TurbineOutput().PressureOutput())
	scheme.massRateSink = sink.SinkPort(scheme.regenerativeGasGenerator.TurbineOutput().MassRateOutput())
	scheme.powerSink = sink.SinkPort(scheme.freeTurbineBlock.PowerOutput())
}
