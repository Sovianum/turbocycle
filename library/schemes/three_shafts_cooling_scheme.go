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

func NewThreeShaftsCoolingScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	middlePressureCascade compose.TurboCascadeNode,
	cooler constructive.CoolerNode,
	gasGenerator compose.GasGeneratorNode,
	middlePressureCompressorPipe constructive.PressureLossNode,
	highPressureTurbinePipe constructive.PressureLossNode,
	middlePressureTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
) ThreeShaftsCoolerScheme {
	return &threeShaftsCoolerScheme{
		gasSource:                    gasSource,
		inletPressureDrop:            inletPressureDrop,
		middlePressureCascade:        middlePressureCascade,
		gasGenerator:                 gasGenerator,
		middlePressureCompressorPipe: middlePressureCompressorPipe,
		cooler:                    cooler,
		highPressureTurbinePipe:   highPressureTurbinePipe,
		middlePressureTurbinePipe: middlePressureTurbinePipe,
		freeTurbineBlock:          freeTurbineBlock,
		breaker: helper.NewComplexCycleBreakNode(
			gases.GetAir(), states.StandardTemperature, states.StandardPressure, 1,
		),
	}
}

type ThreeShaftsCoolerScheme interface {
	Scheme
	DoubleCompressor
}

type threeShaftsCoolerScheme struct {
	gasSource                    source.ComplexGasSourceNode
	inletPressureDrop            constructive.PressureLossNode
	middlePressureCascade        compose.TurboCascadeNode
	gasGenerator                 compose.GasGeneratorNode
	middlePressureCompressorPipe constructive.PressureLossNode
	cooler                       constructive.CoolerNode
	highPressureTurbinePipe      constructive.PressureLossNode
	middlePressureTurbinePipe    constructive.PressureLossNode
	freeTurbineBlock             compose.FreeTurbineBlockNode
	breaker                      helper.ComplexCycleBreakNode

	gasSink         sink.SinkNode
	pressureSink    sink.SinkNode
	temperatureSink sink.SinkNode
	massRateSink    sink.SinkNode
	powerSink       sink.SinkNode
}

func (scheme *threeShaftsCoolerScheme) LowPressureCompressor() constructive.CompressorNode {
	return scheme.middlePressureCascade.Compressor()
}

func (scheme *threeShaftsCoolerScheme) HighPressureCompressor() constructive.CompressorNode {
	return scheme.gasGenerator.TurboCascade().Compressor()
}

func (scheme *threeShaftsCoolerScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.MassRateOutput().GetState().(states.MassRatePortState).MassRate
	return lSpecific * massRateRel
}

func (scheme *threeShaftsCoolerScheme) GetFuelMassRateRel() float64 {
	var burner = scheme.gasGenerator.Burner()
	var massRateRel = burner.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return burner.GetFuelRateRel() * massRateRel
}

func (scheme *threeShaftsCoolerScheme) GetQLower() float64 {
	return scheme.gasGenerator.Burner().Fuel().QLower()
}

func (scheme *threeShaftsCoolerScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()

	return graph.NewNetwork([]graph.Node{
		scheme.gasSource, scheme.inletPressureDrop, scheme.middlePressureCascade, scheme.cooler,
		scheme.middlePressureCompressorPipe, scheme.gasGenerator, scheme.highPressureTurbinePipe,
		scheme.middlePressureTurbinePipe, scheme.freeTurbineBlock,
		scheme.powerSink, scheme.gasSink, scheme.pressureSink, scheme.temperatureSink, scheme.massRateSink,
		scheme.breaker,
	})
}

func (scheme *threeShaftsCoolerScheme) linkPorts() {
	nodes.LinkComplexOutToIn(scheme.gasSource, scheme.inletPressureDrop)
	nodes.LinkComplexOutToIn(scheme.inletPressureDrop, scheme.middlePressureCascade.CompressorComplexGasInput())

	nodes.LinkComplexOutToIn(scheme.middlePressureCascade.CompressorComplexGasOutput(), scheme.middlePressureCompressorPipe)
	nodes.LinkComplexOutToIn(scheme.middlePressureCompressorPipe, scheme.cooler)
	nodes.LinkComplexOutToIn(scheme.cooler, scheme.breaker)
	nodes.LinkComplexOutToIn(scheme.breaker, scheme.gasGenerator)
	nodes.LinkComplexOutToIn(scheme.gasGenerator, scheme.highPressureTurbinePipe)
	nodes.LinkComplexOutToIn(scheme.highPressureTurbinePipe, scheme.middlePressureCascade.TurbineComplexGasInput())
	nodes.LinkComplexOutToIn(scheme.middlePressureCascade.TurbineComplexGasOutput(), scheme.middlePressureTurbinePipe)
	nodes.LinkComplexOutToIn(scheme.middlePressureTurbinePipe, scheme.freeTurbineBlock)

	scheme.powerSink = sink.SinkPort(scheme.freeTurbineBlock.PowerOutput())
	scheme.gasSink = sink.SinkPort(scheme.freeTurbineBlock.GasOutput())
	scheme.pressureSink = sink.SinkPort(scheme.freeTurbineBlock.PressureOutput())
	scheme.temperatureSink = sink.SinkPort(scheme.freeTurbineBlock.TemperatureOutput())
	scheme.massRateSink = sink.SinkPort(scheme.freeTurbineBlock.MassRateOutput())
}
