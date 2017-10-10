package schemes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/nodes/source"
	"github.com/Sovianum/turbocycle/impl/states"
)

func NewThreeShaftsRegeneratorScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	middlePressureCascade compose.TurboCascadeNode,
	regenerativeGasGenerator compose.RegenerativeGasGeneratorNode,
	middlePressureCompressorPipe constructive.PressureLossNode,
	highPressureTurbinePipe constructive.PressureLossNode,
	middlePressureTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
) ThreeShaftsRegeneratorScheme {
	return &threeShaftsRegeneratorScheme{
		gasSource:                    gasSource,
		inletPressureDrop:            inletPressureDrop,
		middlePressureCascade:        middlePressureCascade,
		regenerativeGasGenerator:     regenerativeGasGenerator,
		middlePressureCompressorPipe: middlePressureCompressorPipe,
		highPressureTurbinePipe:      highPressureTurbinePipe,
		middlePressureTurbinePipe:    middlePressureTurbinePipe,
		freeTurbineBlock:             freeTurbineBlock,
		gasSink:                      sink.NewComplexGasSinkNode(),
		powerSink:                    sink.NewPowerSinkNode(),
		breaker1:                     helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
		breaker2:                     helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
	}
}

type ThreeShaftsRegeneratorScheme interface {
	Scheme
	InitGasGeneratorCompressor(state states.ComplexGasPortState)
	InitGasGeneratorHeatExchanger(state states.ComplexGasPortState)
}

type threeShaftsRegeneratorScheme struct {
	gasSource                    source.ComplexGasSourceNode
	inletPressureDrop            constructive.PressureLossNode
	middlePressureCascade        compose.TurboCascadeNode
	regenerativeGasGenerator     compose.RegenerativeGasGeneratorNode
	middlePressureCompressorPipe constructive.PressureLossNode
	highPressureTurbinePipe      constructive.PressureLossNode
	middlePressureTurbinePipe    constructive.PressureLossNode
	freeTurbineBlock             compose.FreeTurbineBlockNode
	gasSink                      sink.ComplexGasSinkNode
	powerSink                    nodes.PowerSink
	breaker1                     helper.CycleBreakNode
	breaker2                     helper.CycleBreakNode
}

func (scheme *threeShaftsRegeneratorScheme) InitGasGeneratorCompressor(state states.ComplexGasPortState) {
	scheme.breaker1.DataSourcePort().SetState(state)
}

func (scheme *threeShaftsRegeneratorScheme) InitGasGeneratorHeatExchanger(state states.ComplexGasPortState) {
	scheme.breaker2.DataSourcePort().SetState(state)
}

func (scheme *threeShaftsRegeneratorScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return lSpecific * massRateRel
}

func (scheme *threeShaftsRegeneratorScheme) GetFuelMassRateRel() float64 {
	var massRateRel = scheme.regenerativeGasGenerator.Burner().ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return scheme.regenerativeGasGenerator.Burner().GetFuelRateRel() * massRateRel
}

func (scheme *threeShaftsRegeneratorScheme) GetQLower() float64 {
	return scheme.regenerativeGasGenerator.Burner().Fuel().QLower()
}


func (scheme *threeShaftsRegeneratorScheme) GetNetwork() core.Network {
	var nodeMap = make(map[string]core.Node)
	nodeMap[inputGasSourceName] = scheme.gasSource
	nodeMap[inletPressureDropName] = scheme.inletPressureDrop
	nodeMap[middlePressureCascadeName] = scheme.middlePressureCascade
	nodeMap[middlePressureCompressorPipeName] = scheme.middlePressureCompressorPipe
	nodeMap[regenerativeGasGeneratorName] = scheme.regenerativeGasGenerator
	nodeMap[highPressureTurbinePipeName] = scheme.highPressureTurbinePipe
	nodeMap[middlePressureTurbinePipeName] = scheme.middlePressureTurbinePipe
	nodeMap[freeTurbineBlockName] = scheme.freeTurbineBlock
	nodeMap[outputGasSinkName] = scheme.gasSink
	nodeMap["breaker1"] = scheme.breaker1
	nodeMap["breaker2"] = scheme.breaker2

	core.Link(scheme.gasSource.ComplexGasOutput(), scheme.inletPressureDrop.ComplexGasInput())
	core.Link(scheme.inletPressureDrop.ComplexGasOutput(), scheme.middlePressureCascade.CompressorComplexGasInput())

	core.Link(scheme.middlePressureCascade.CompressorComplexGasOutput(), scheme.middlePressureCompressorPipe.ComplexGasInput())
	core.Link(scheme.middlePressureCompressorPipe.ComplexGasOutput(), scheme.breaker1.DataSourcePort())
	core.Link(scheme.breaker1.UpdatePort(), scheme.regenerativeGasGenerator.ComplexGasInput())
	core.Link(scheme.regenerativeGasGenerator.ComplexGasOutput(), scheme.highPressureTurbinePipe.ComplexGasInput())
	core.Link(scheme.highPressureTurbinePipe.ComplexGasOutput(), scheme.middlePressureCascade.TurbineComplexGasInput())
	core.Link(scheme.middlePressureCascade.TurbineComplexGasOutput(), scheme.middlePressureTurbinePipe.ComplexGasInput())
	core.Link(scheme.middlePressureTurbinePipe.ComplexGasOutput(), scheme.freeTurbineBlock.ComplexGasInput())
	core.Link(scheme.freeTurbineBlock.PowerOutput(), scheme.powerSink.PowerInput())
	core.Link(scheme.freeTurbineBlock.ComplexGasOutput(), scheme.breaker2.DataSourcePort())
	core.Link(scheme.breaker2.UpdatePort(), scheme.regenerativeGasGenerator.HeatExchangerHotInput())
	core.Link(scheme.regenerativeGasGenerator.HeatExchangerHotOutput(), scheme.gasSink.ComplexGasInput())

	return core.NewNetwork(nodeMap)
}
