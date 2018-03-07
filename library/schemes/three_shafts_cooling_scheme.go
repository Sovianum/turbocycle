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
	LPCPipe() constructive.PressureLossNode
	HPCPipe() constructive.PressureLossNode
	MainBurner() constructive.BurnerNode
	HPT() constructive.StaticTurbineNode
	LPT() constructive.StaticTurbineNode
	FT() constructive.StaticTurbineNode
	FTBlock() compose.FreeTurbineBlockNode
	LPTPipe() constructive.PressureLossNode
	HPTPipe() constructive.PressureLossNode
	GasGenerator() compose.GasGeneratorNode
	MiddlePressureCascade() compose.TurboCascadeNode
	InletPressureDrop() constructive.PressureLossNode
	GasSource() source.ComplexGasSourceNode
	Cooler() constructive.CoolerNode
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

func (scheme *threeShaftsCoolerScheme) Cooler() constructive.CoolerNode {
	return scheme.cooler
}

func (scheme *threeShaftsCoolerScheme) LPCPipe() constructive.PressureLossNode {
	return scheme.middlePressureCompressorPipe
}

func (scheme *threeShaftsCoolerScheme) HPCPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(1) // todo make smth more precise
}

func (scheme *threeShaftsCoolerScheme) MainBurner() constructive.BurnerNode {
	return scheme.gasGenerator.Burner()
}

func (scheme *threeShaftsCoolerScheme) HPT() constructive.StaticTurbineNode {
	return scheme.gasGenerator.TurboCascade().Turbine()
}

func (scheme *threeShaftsCoolerScheme) LPT() constructive.StaticTurbineNode {
	return scheme.middlePressureCascade.Turbine()
}

func (scheme *threeShaftsCoolerScheme) FT() constructive.StaticTurbineNode {
	return scheme.freeTurbineBlock.FreeTurbine()
}

func (scheme *threeShaftsCoolerScheme) FTBlock() compose.FreeTurbineBlockNode {
	return scheme.freeTurbineBlock
}

func (scheme *threeShaftsCoolerScheme) LPTPipe() constructive.PressureLossNode {
	return scheme.middlePressureTurbinePipe
}

func (scheme *threeShaftsCoolerScheme) HPTPipe() constructive.PressureLossNode {
	return scheme.highPressureTurbinePipe
}

func (scheme *threeShaftsCoolerScheme) GasGenerator() compose.GasGeneratorNode {
	return scheme.gasGenerator
}

func (scheme *threeShaftsCoolerScheme) MiddlePressureCascade() compose.TurboCascadeNode {
	return scheme.middlePressureCascade
}

func (scheme *threeShaftsCoolerScheme) InletPressureDrop() constructive.PressureLossNode {
	return scheme.inletPressureDrop
}

func (scheme *threeShaftsCoolerScheme) GasSource() source.ComplexGasSourceNode {
	return scheme.gasSource
}

func (scheme *threeShaftsCoolerScheme) LPC() constructive.CompressorNode {
	return scheme.middlePressureCascade.Compressor()
}

func (scheme *threeShaftsCoolerScheme) HPC() constructive.CompressorNode {
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
	return burner.FuelRateRel() * massRateRel
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
