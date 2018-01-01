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
		gasSink:                   sink.NewComplexGasSinkNode(),
		powerSink:                 sink.NewPowerSinkNode(),
		breaker:                   helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
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
	gasSink                      sink.ComplexGasSinkNode
	powerSink                    nodes.PowerSink
	breaker                      helper.CycleBreakNode
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
	var massRateRel = turbine.ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return lSpecific * massRateRel
}

func (scheme *threeShaftsCoolerScheme) GetFuelMassRateRel() float64 {
	var massRateRel = scheme.gasGenerator.Burner().ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return scheme.gasGenerator.Burner().GetFuelRateRel() * massRateRel
}

func (scheme *threeShaftsCoolerScheme) GetQLower() float64 {
	return scheme.gasGenerator.Burner().Fuel().QLower()
}

func (scheme *threeShaftsCoolerScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()

	return graph.NewNetwork([]graph.Node{
		scheme.gasSource, scheme.inletPressureDrop, scheme.middlePressureCascade, scheme.cooler,
		scheme.middlePressureCompressorPipe, scheme.gasGenerator, scheme.highPressureTurbinePipe,
		scheme.middlePressureTurbinePipe, scheme.freeTurbineBlock, scheme.gasSink, scheme.breaker,
	})
}

func (scheme *threeShaftsCoolerScheme) linkPorts() {
	graph.Link(scheme.gasSource.ComplexGasOutput(), scheme.inletPressureDrop.ComplexGasInput())
	graph.Link(scheme.inletPressureDrop.ComplexGasOutput(), scheme.middlePressureCascade.CompressorComplexGasInput())

	graph.Link(scheme.middlePressureCascade.CompressorComplexGasOutput(), scheme.middlePressureCompressorPipe.ComplexGasInput())
	graph.Link(scheme.middlePressureCompressorPipe.ComplexGasOutput(), scheme.cooler.ComplexGasInput())
	graph.Link(scheme.cooler.ComplexGasOutput(), scheme.breaker.DataSourcePort())
	graph.Link(scheme.breaker.UpdatePort(), scheme.gasGenerator.ComplexGasInput())
	graph.Link(scheme.gasGenerator.ComplexGasOutput(), scheme.highPressureTurbinePipe.ComplexGasInput())
	graph.Link(scheme.highPressureTurbinePipe.ComplexGasOutput(), scheme.middlePressureCascade.TurbineComplexGasInput())
	graph.Link(scheme.middlePressureCascade.TurbineComplexGasOutput(), scheme.middlePressureTurbinePipe.ComplexGasInput())
	graph.Link(scheme.middlePressureTurbinePipe.ComplexGasOutput(), scheme.freeTurbineBlock.ComplexGasInput())
	graph.Link(scheme.freeTurbineBlock.PowerOutput(), scheme.powerSink.PowerInput())
	graph.Link(scheme.freeTurbineBlock.ComplexGasOutput(), scheme.gasSink.ComplexGasInput())
}
