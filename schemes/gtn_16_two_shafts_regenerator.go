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

func NewGtn16TwoShaftsRegenerator(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	turboCascade compose.TurboCascadeNode,
	burner constructive.BurnerNode,
	compressorTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
	regenerator constructive.RegeneratorNode,

) Gtn16TwoShaftsRegenerator {
	return &gtn16TwoShaftsRegenerator{
		gasSource:             gasSource,
		inletPressureDrop:     inletPressureDrop,
		turboCascade:          turboCascade,
		burner:                burner,
		compressorTurbinePipe: compressorTurbinePipe,
		freeTurbineBlock:      freeTurbineBlock,
		regenerator:           regenerator,
		gasSink:               sink.NewComplexGasSinkNode(),
		powerSink:             sink.NewPowerSinkNode(),
		breaker1:              helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
		breaker2:              helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
	}
}

type Gtn16TwoShaftsRegenerator interface {
	Scheme
}

type gtn16TwoShaftsRegenerator struct {
	gasSource             source.ComplexGasSourceNode
	inletPressureDrop     constructive.PressureLossNode
	turboCascade          compose.TurboCascadeNode
	burner                constructive.BurnerNode
	compressorTurbinePipe constructive.PressureLossNode
	freeTurbineBlock      compose.FreeTurbineBlockNode
	regenerator           constructive.RegeneratorNode
	gasSink               sink.ComplexGasSinkNode
	powerSink             nodes.PowerSink
	breaker1              helper.CycleBreakNode
	breaker2              helper.CycleBreakNode
}

func (scheme *gtn16TwoShaftsRegenerator) GetNetwork() core.Network {
	var nodeMap = make(map[string]core.Node)
	nodeMap[inputGasSourceName] = scheme.gasSource
	nodeMap[inletPressureDropName] = scheme.inletPressureDrop
	nodeMap[turboCascadeName] = scheme.turboCascade
	nodeMap[regeneratorName] = scheme.regenerator
	nodeMap[burnerName] = scheme.burner
	nodeMap[compressorTurbinePipeName] = scheme.compressorTurbinePipe
	nodeMap[freeTurbineBlockName] = scheme.freeTurbineBlock
	nodeMap[outputGasSinkName] = scheme.gasSink
	nodeMap["breaker1"] = scheme.breaker1
	nodeMap["breaker2"] = scheme.breaker2

	core.Link(scheme.gasSource.ComplexGasOutput(), scheme.inletPressureDrop.ComplexGasInput())
	core.Link(scheme.inletPressureDrop.ComplexGasOutput(), scheme.turboCascade.CompressorComplexGasInput())
	core.Link(scheme.turboCascade.CompressorComplexGasOutput(), scheme.regenerator.ColdInput())
	core.Link(scheme.regenerator.ColdOutput(), scheme.breaker1.DataSourcePort())
	core.Link(scheme.breaker1.UpdatePort(), scheme.burner.ComplexGasInput())
	core.Link(scheme.burner.ComplexGasOutput(), scheme.turboCascade.TurbineComplexGasInput())
	core.Link(scheme.turboCascade.TurbineComplexGasOutput(), scheme.compressorTurbinePipe.ComplexGasInput())
	core.Link(scheme.compressorTurbinePipe.ComplexGasOutput(), scheme.freeTurbineBlock.ComplexGasInput())
	core.Link(scheme.freeTurbineBlock.ComplexGasOutput(), scheme.breaker2.DataSourcePort())
	core.Link(scheme.breaker2.UpdatePort(), scheme.regenerator.HotInput())
	core.Link(scheme.regenerator.HotOutput(), scheme.gasSink.ComplexGasInput())
	core.Link(scheme.freeTurbineBlock.PowerOutput(), scheme.powerSink.PowerInput())

	return core.NewNetwork(nodeMap)
}
