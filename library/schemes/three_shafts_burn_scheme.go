package schemes

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

func NewThreeShaftsBurnScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	middlePressureCascade compose.TurboCascadeNode,
	gasGenerator compose.GasGeneratorNode,
	middlePressureCompressorPipe constructive.PressureLossNode,
	highPressureTurbinePipe constructive.PressureLossNode,
	middlePressureTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
	midBurner constructive.BurnerNode,
) ThreeShaftsBurnScheme {
	return &threeShaftsBurnScheme{
		threeShaftsScheme: &threeShaftsScheme{
			gasSource:                    gasSource,
			inletPressureDrop:            inletPressureDrop,
			middlePressureCascade:        middlePressureCascade,
			gasGenerator:                 gasGenerator,
			middlePressureCompressorPipe: middlePressureCompressorPipe,
			highPressureTurbinePipe:      highPressureTurbinePipe,
			middlePressureTurbinePipe:    middlePressureTurbinePipe,
			freeTurbineBlock:             freeTurbineBlock,
			breaker: helper.NewComplexCycleBreakNode(
				gases.GetAir(), states.StandardTemperature, states.StandardPressure, 1,
			),
		},
		midBurner: midBurner,
	}
}

type ThreeShaftsBurnScheme interface {
	ThreeShaftsScheme
	MidBurner() constructive.BurnerNode
}

type threeShaftsBurnScheme struct {
	*threeShaftsScheme
	midBurner constructive.BurnerNode
}

func (scheme *threeShaftsBurnScheme) GetFuelMassRateRel() float64 {
	mainBurner := scheme.gasGenerator.Burner()
	mainMassRateRel := mainBurner.MassRateInput().GetState().(states.MassRatePortState).MassRate
	mainFuelRate := mainBurner.FuelRateRel() * mainMassRateRel

	extraBurner := scheme.midBurner
	extraMassRateRel := extraBurner.MassRateInput().GetState().(states.MassRatePortState).MassRate
	extraFuelRate := extraBurner.FuelRateRel() * extraMassRateRel

	return mainFuelRate + extraFuelRate
}

func (scheme *threeShaftsBurnScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()
	baseNodes := scheme.threeShaftsScheme.nodes()
	return graph.NewNetwork(append(baseNodes, scheme.midBurner))
}

func (scheme *threeShaftsBurnScheme) MidBurner() constructive.BurnerNode {
	return scheme.midBurner
}

func (scheme *threeShaftsBurnScheme) linkPorts() {
	scheme.threeShaftsScheme.linkPorts()
	nodes.LinkComplexOutToIn(scheme.highPressureTurbinePipe, scheme.midBurner)
	nodes.LinkComplexOutToIn(scheme.midBurner, scheme.middlePressureCascade.TurbineComplexGasInput())
}
