package compose

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
)

type TurboCascadeNode interface {
	graph.Node
	CompressorComplexGasInput() nodes.ComplexGasSink
	CompressorComplexGasOutput() nodes.ComplexGasSource
	TurbineComplexGasInput() nodes.ComplexGasSink
	TurbineComplexGasOutput() nodes.ComplexGasSource
	Compressor() constructive.CompressorNode
	Turbine() constructive.StaticTurbineNode
	Transmission() constructive.TransmissionNode
	ProcessCompressor() error
	ProcessTurbine() error
}

func NewTurboCascadeNode(
	compressorEtaAd, piStag float64,
	etaT, lambdaOut float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(constructive.TurbineNode) float64,
	etaM float64,
	precision float64,
) TurboCascadeNode {
	var result = &turboCascadeNode{
		compressor: constructive.NewCompressorNode(compressorEtaAd, piStag, precision),
		turbine: constructive.NewBlockedTurbineNode(
			etaT, lambdaOut, precision,
			leakMassRateFunc, coolMasRateRel, inflowMassRateRel,
		),
		transmission: constructive.NewTransmissionNode(etaM),
	}

	result.linkPorts()

	graph.AttachAllWithTags(
		result,
		[]*graph.Port{
			&result.compressorGasInput, &result.compressorTemperatureInput,
			&result.compressorPressureInput, &result.compressorMassRateInput,

			&result.compressorGasOutput, &result.compressorTemperatureOutput,
			&result.compressorPressureOutput, &result.compressorMassRateOutput,

			&result.turbineGasInput, &result.turbineTemperatureInput,
			&result.turbinePressureInput, &result.turbineMassRateInput,

			&result.turbineGasOutput, &result.turbineTemperatureOutput,
			&result.turbinePressureOutput, &result.turbineMassRateOutput,
		},
		[]string{
			"compGasInput", "compTempInput", "compPressureInput", "compMassRateInput",
			"compGasOutput", "compTempOutput", "compPressureOutput", "compMassRateOutput",
			"turbGasInput", "turbTempInput", "turbPressureInput", "turbMassRateInput",
			"turbGasOutput", "turbTempOutput", "turbPressureOutput", "turbMassRateOutput",
		},
	)

	return result
}

type turboCascadeNode struct {
	graph.BaseNode

	compressorGasInput         graph.Port
	compressorTemperatureInput graph.Port
	compressorPressureInput    graph.Port
	compressorMassRateInput    graph.Port

	compressorGasOutput         graph.Port
	compressorTemperatureOutput graph.Port
	compressorPressureOutput    graph.Port
	compressorMassRateOutput    graph.Port

	turbineGasInput         graph.Port
	turbineTemperatureInput graph.Port
	turbinePressureInput    graph.Port
	turbineMassRateInput    graph.Port

	turbineGasOutput         graph.Port
	turbineTemperatureOutput graph.Port
	turbinePressureOutput    graph.Port
	turbineMassRateOutput    graph.Port

	compressor   constructive.CompressorNode
	turbine      constructive.BlockedTurbineNode
	transmission constructive.TransmissionNode
}

func (node *turboCascadeNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "TurboCascade")
}

func (node *turboCascadeNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.compressorGasInput,
		node.compressorTemperatureInput,
		node.compressorPressureInput,
		node.compressorMassRateInput,

		node.compressorGasOutput,
		node.compressorTemperatureOutput,
		node.compressorPressureOutput,
		node.compressorMassRateOutput,

		node.turbineGasInput,
		node.turbineTemperatureInput,
		node.turbinePressureInput,
		node.turbineMassRateInput,

		node.turbineGasOutput,
		node.turbineTemperatureOutput,
		node.turbinePressureOutput,
		node.turbineMassRateOutput,
	}
}

func (node *turboCascadeNode) GetRequirePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.compressorGasInput,
		node.compressorTemperatureInput,
		node.compressorPressureInput,
		node.compressorMassRateInput,

		node.turbineGasInput,
		node.turbineTemperatureInput,
		node.turbinePressureInput,
		node.turbineMassRateInput,
	}, nil
}

func (node *turboCascadeNode) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.compressorGasOutput,
		node.compressorTemperatureOutput,
		node.compressorPressureOutput,
		node.compressorMassRateOutput,

		node.turbineGasOutput,
		node.turbineTemperatureOutput,
		node.turbinePressureOutput,
		node.turbineMassRateOutput,
	}, nil
}

func (node *turboCascadeNode) Compressor() constructive.CompressorNode {
	return node.compressor
}

func (node *turboCascadeNode) Turbine() constructive.StaticTurbineNode {
	return node.turbine
}

func (node *turboCascadeNode) Transmission() constructive.TransmissionNode {
	return node.transmission
}

func (node *turboCascadeNode) ProcessCompressor() error {
	node.readInput()
	if err := node.compressor.Process(); err != nil {
		return err
	}
	if err := node.transmission.Process(); err != nil {
		return err
	}
	node.writeOutput()

	return nil
}

func (node *turboCascadeNode) ProcessTurbine() error {
	node.readInput()
	if err := node.turbine.Process(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *turboCascadeNode) Process() error {
	node.readInput()
	if err := node.compressor.Process(); err != nil {
		return err
	}
	if err := node.transmission.Process(); err != nil {
		return err
	}
	if err := node.turbine.Process(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *turboCascadeNode) TurbineComplexGasOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.turbineGasOutput, node.turbineTemperatureOutput,
		node.turbinePressureOutput, node.turbineMassRateOutput,
	)
}

func (node *turboCascadeNode) TurbineComplexGasInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.turbineGasInput, node.turbineTemperatureInput,
		node.turbinePressureInput, node.turbineMassRateInput,
	)
}

func (node *turboCascadeNode) CompressorComplexGasOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.compressorGasOutput, node.compressorTemperatureOutput,
		node.compressorPressureOutput, node.compressorMassRateOutput,
	)
}

func (node *turboCascadeNode) CompressorComplexGasInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.compressorGasInput, node.compressorTemperatureInput,
		node.compressorPressureInput, node.compressorMassRateInput,
	)
}

func (node *turboCascadeNode) linkPorts() {
	graph.Link(node.compressor.PowerOutput(), node.transmission.PowerInput())
	graph.Link(node.transmission.PowerOutput(), node.turbine.PowerInput())
}

func (node *turboCascadeNode) readInput() {
	graph.SetAll(
		[]graph.PortState{
			node.compressorGasInput.GetState(), node.compressorTemperatureInput.GetState(),
			node.compressorPressureInput.GetState(), node.compressorMassRateInput.GetState(),
		},
		[]graph.Port{
			node.compressor.GasInput(), node.compressor.TemperatureInput(),
			node.compressor.PressureInput(), node.compressor.MassRateInput(),
		},
	)

	graph.SetAll(
		[]graph.PortState{
			node.turbineGasInput.GetState(), node.turbineTemperatureInput.GetState(),
			node.turbinePressureInput.GetState(), node.turbineMassRateInput.GetState(),
		},
		[]graph.Port{
			node.turbine.GasInput(), node.turbine.TemperatureInput(),
			node.turbine.PressureInput(), node.turbine.MassRateInput(),
		},
	)
}

func (node *turboCascadeNode) writeOutput() {
	graph.SetAll(
		[]graph.PortState{
			node.compressor.GasOutput().GetState(), node.compressor.TemperatureOutput().GetState(),
			node.compressor.PressureOutput().GetState(), node.compressor.MassRateOutput().GetState(),
		},
		[]graph.Port{
			node.compressorGasOutput, node.compressorTemperatureOutput,
			node.compressorPressureOutput, node.compressorMassRateOutput,
		},
	)

	graph.SetAll(
		[]graph.PortState{
			node.turbine.GasOutput().GetState(), node.turbine.TemperatureOutput().GetState(),
			node.turbine.PressureOutput().GetState(), node.turbine.MassRateOutput().GetState(),
		},
		[]graph.Port{
			node.turbineGasOutput, node.turbineTemperatureOutput,
			node.turbinePressureOutput, node.turbineMassRateOutput,
		},
	)
}
