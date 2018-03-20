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

func NewThreeShaftsSubCompressScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	middlePressureCascade compose.TurboCascadeNode,
	gasGenerator compose.GasGeneratorNode,
	middlePressureCompressorPipe constructive.PressureLossNode,
	highPressureTurbinePipe constructive.PressureLossNode,
	middlePressureTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
	gasSplitter constructive.GasSplitter,
	gasCombiner constructive.GasCombiner,
	subCompressor constructive.CompressorNode,
	subCooler constructive.CoolerNode,
) ThreeShaftsSubCompressScheme {
	return &threeShaftsSubCompressScheme{
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
		gasSplitter:     gasSplitter,
		gasCombiner:     gasCombiner,
		subCompressor:   subCompressor,
		powerAggregator: constructive.NewPowerAggregator(),
		subCooler:       subCooler,
	}
}

type ThreeShaftsSubCompressScheme interface {
	ThreeShaftsScheme
	SubCompressor() constructive.CompressorNode
	SubCooler() constructive.CoolerNode
	GasSplitter() constructive.GasSplitter
	GasCombiner() constructive.GasCombiner
	PowerAggregator() constructive.PowerAggregator
}

type threeShaftsSubCompressScheme struct {
	*threeShaftsScheme

	gasSplitter     constructive.GasSplitter
	gasCombiner     constructive.GasCombiner
	subCompressor   constructive.CompressorNode
	subCooler       constructive.CoolerNode
	powerAggregator constructive.PowerAggregator
}

func (scheme *threeShaftsSubCompressScheme) SubCompressor() constructive.CompressorNode {
	return scheme.subCompressor
}

func (scheme *threeShaftsSubCompressScheme) SubCooler() constructive.CoolerNode {
	return scheme.subCooler
}

func (scheme *threeShaftsSubCompressScheme) GasSplitter() constructive.GasSplitter {
	return scheme.gasSplitter
}

func (scheme *threeShaftsSubCompressScheme) GasCombiner() constructive.GasCombiner {
	return scheme.gasCombiner
}

func (scheme *threeShaftsSubCompressScheme) PowerAggregator() constructive.PowerAggregator {
	return scheme.powerAggregator
}

func (scheme *threeShaftsSubCompressScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()
	return graph.NewNetwork(scheme.nodes())
}

func (scheme *threeShaftsSubCompressScheme) nodes() []graph.Node {
	disassembledNodes := []graph.Node{ // nodes are disassembled to prevent its processing as a whole
		scheme.gasSource,
		scheme.inletPressureDrop,
		scheme.LPC(),
		scheme.LPCPipe(),
		scheme.HPC(),
		scheme.MainBurner(),
		scheme.HPT(),
		scheme.HPTPipe(),
		scheme.LPT(),
		scheme.LPTPipe(),
		scheme.gasGenerator.TurboCascade().Transmission(),
		scheme.middlePressureCascade.Transmission(),
		scheme.freeTurbineBlock,
	}

	return append(
		disassembledNodes,
		scheme.gasSplitter,
		scheme.gasCombiner,
		scheme.subCompressor,
		scheme.powerAggregator,
		scheme.subCooler,
	)
}

func (scheme *threeShaftsSubCompressScheme) linkPorts() {
	scheme.HPC().SetName("hpc")
	scheme.HPT().SetName("hpt")
	scheme.LPC().SetName("lpc")
	scheme.LPT().SetName("lpt")
	scheme.FT().SetName("ft")
	scheme.HPTPipe().SetName("hptPipe")
	scheme.LPCPipe().SetName("lpcPipe")
	scheme.LPTPipe().SetName("lptPipe")
	scheme.inletPressureDrop.SetName("inletPressureDrop")
	scheme.subCompressor.SetName("subCompressor")

	scheme.linkDisassemble()

	// make extra gas path. Taken gas is passed from hpc around hpt directly to mpt
	nodes.LinkComplexOutToIn(scheme.HPC(), scheme.gasSplitter.Input())
	nodes.LinkComplexOutToIn(scheme.gasSplitter.MainOutput(), scheme.MainBurner())
	nodes.LinkComplexOutToIn(scheme.gasSplitter.ExtraOutput(), scheme.subCooler)
	nodes.LinkComplexOutToIn(scheme.subCooler, scheme.subCompressor)
	nodes.LinkComplexOutToIn(scheme.subCompressor, scheme.gasCombiner.ExtraInput())
	// repair main gas path: hpt pipe connected to gas combiner
	nodes.LinkComplexOutToIn(scheme.HPTPipe(), scheme.gasCombiner.MainInput())
	nodes.LinkComplexOutToIn(scheme.gasCombiner.Output(), scheme.LPT())

	// make extra power path. Power demand from subcompressor is compensated with
	// middle pressure turbine
	graph.Link(
		graph.NewWeakPort(scheme.subCompressor.MassRateInput()),
		scheme.powerAggregator.ExtraMassRateInput(),
	)
	graph.Link(
		scheme.subCompressor.PowerOutput(),
		scheme.powerAggregator.ExtraPowerInput(),
	)
	graph.Link(
		graph.NewWeakPort(scheme.LPC().MassRateInput()),
		scheme.powerAggregator.MainMassRateInput(),
	)
	graph.Link(
		scheme.MiddlePressureCascade().Transmission().PowerOutput(),
		scheme.powerAggregator.MainPowerInput(),
	)
	graph.Link(
		scheme.powerAggregator.PowerOutput(),
		scheme.LPT().(constructive.BlockedTurbineNode).PowerInput(),
	)
}

// this method disassembles all complex nodes nodes of base scheme
// to prevent cyclic dependencies
func (scheme *threeShaftsSubCompressScheme) linkDisassemble() {
	nodes.StartLink(scheme.gasSource).
		Link(scheme.inletPressureDrop).
		Link(scheme.LPC()).
		Link(scheme.LPCPipe()).
		Link(scheme.HPC()).
		Link(scheme.MainBurner()).
		Link(scheme.HPT()).
		Link(scheme.HPTPipe()).
		Link(scheme.LPT()).
		Link(scheme.LPTPipe()).
		FinishLink(scheme.freeTurbineBlock)

	graph.Link(
		scheme.middlePressureCascade.Transmission().PowerOutput(),
		scheme.LPT().(constructive.BlockedTurbineNode).PowerInput(),
	)
	graph.Link(
		scheme.gasGenerator.TurboCascade().Transmission().PowerOutput(),
		scheme.HPT().(constructive.BlockedTurbineNode).PowerInput(),
	)
	sink.SinkAll(
		scheme.HPT().PowerOutput(), scheme.LPT().PowerOutput(),
		scheme.freeTurbineBlock.GasOutput(),
		scheme.freeTurbineBlock.TemperatureOutput(),
		scheme.freeTurbineBlock.PressureOutput(),
		scheme.freeTurbineBlock.MassRateOutput(),
		scheme.freeTurbineBlock.PowerOutput(),
	)
}
