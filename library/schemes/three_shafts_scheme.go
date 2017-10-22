package schemes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

func NewThreeShaftsScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	middlePressureCascade compose.TurboCascadeNode,
	gasGenerator compose.GasGeneratorNode,
	middlePressureCompressorPipe constructive.PressureLossNode,
	highPressureTurbinePipe constructive.PressureLossNode,
	middlePressureTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
) ThreeShaftsScheme {
	return &threeShaftsScheme{
		gasSource:                    gasSource,
		inletPressureDrop:            inletPressureDrop,
		middlePressureCascade:        middlePressureCascade,
		gasGenerator:                 gasGenerator,
		middlePressureCompressorPipe: middlePressureCompressorPipe,
		highPressureTurbinePipe:      highPressureTurbinePipe,
		middlePressureTurbinePipe:    middlePressureTurbinePipe,
		freeTurbineBlock:             freeTurbineBlock,
		gasSink:                      sink.NewComplexGasSinkNode(),
		powerSink:                    sink.NewPowerSinkNode(),
		breaker:                      helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
	}
}

type ThreeShaftsScheme interface {
	Scheme
	DoubleCompressor
}

type threeShaftsScheme struct {
	gasSource                    source.ComplexGasSourceNode
	inletPressureDrop            constructive.PressureLossNode
	middlePressureCascade        compose.TurboCascadeNode
	gasGenerator                 compose.GasGeneratorNode
	middlePressureCompressorPipe constructive.PressureLossNode
	highPressureTurbinePipe      constructive.PressureLossNode
	middlePressureTurbinePipe    constructive.PressureLossNode
	freeTurbineBlock             compose.FreeTurbineBlockNode
	gasSink                      sink.ComplexGasSinkNode
	powerSink                    nodes.PowerSink
	breaker                      helper.CycleBreakNode
}

func (scheme *threeShaftsScheme) LowPressureCompressor() constructive.CompressorNode {
	return scheme.middlePressureCascade.Compressor()
}

func (scheme *threeShaftsScheme) HighPressureCompressor() constructive.CompressorNode {
	return scheme.gasGenerator.TurboCascade().Compressor()
}

func (scheme *threeShaftsScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return lSpecific * massRateRel
}

func (scheme *threeShaftsScheme) GetFuelMassRateRel() float64 {
	var massRateRel = scheme.gasGenerator.Burner().ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return scheme.gasGenerator.Burner().GetFuelRateRel() * massRateRel
}

func (scheme *threeShaftsScheme) GetQLower() float64 {
	return scheme.gasGenerator.Burner().Fuel().QLower()
}

func (scheme *threeShaftsScheme) GetNetwork() core.Network {
	var nodeMap = make(map[string]core.Node)
	nodeMap[inputGasSourceName] = scheme.gasSource
	nodeMap[inletPressureDropName] = scheme.inletPressureDrop
	nodeMap[middlePressureCascadeName] = scheme.middlePressureCascade
	nodeMap[middlePressureCompressorPipeName] = scheme.middlePressureCompressorPipe
	nodeMap[gasGeneratorName] = scheme.gasGenerator
	nodeMap[highPressureTurbinePipeName] = scheme.highPressureTurbinePipe
	nodeMap[middlePressureTurbinePipeName] = scheme.middlePressureTurbinePipe
	nodeMap[freeTurbineBlockName] = scheme.freeTurbineBlock
	nodeMap[outputGasSinkName] = scheme.gasSink
	nodeMap["breaker"] = scheme.breaker

	core.Link(scheme.gasSource.ComplexGasOutput(), scheme.inletPressureDrop.ComplexGasInput())
	core.Link(scheme.inletPressureDrop.ComplexGasOutput(), scheme.middlePressureCascade.CompressorComplexGasInput())

	core.Link(scheme.middlePressureCascade.CompressorComplexGasOutput(), scheme.middlePressureCompressorPipe.ComplexGasInput())
	core.Link(scheme.middlePressureCompressorPipe.ComplexGasOutput(), scheme.breaker.DataSourcePort())
	core.Link(scheme.breaker.UpdatePort(), scheme.gasGenerator.ComplexGasInput())
	core.Link(scheme.gasGenerator.ComplexGasOutput(), scheme.highPressureTurbinePipe.ComplexGasInput())
	core.Link(scheme.highPressureTurbinePipe.ComplexGasOutput(), scheme.middlePressureCascade.TurbineComplexGasInput())
	core.Link(scheme.middlePressureCascade.TurbineComplexGasOutput(), scheme.middlePressureTurbinePipe.ComplexGasInput())
	core.Link(scheme.middlePressureTurbinePipe.ComplexGasOutput(), scheme.freeTurbineBlock.ComplexGasInput())
	core.Link(scheme.freeTurbineBlock.PowerOutput(), scheme.powerSink.PowerInput())
	core.Link(scheme.freeTurbineBlock.ComplexGasOutput(), scheme.gasSink.ComplexGasInput())

	return core.NewNetwork(nodeMap)
}
