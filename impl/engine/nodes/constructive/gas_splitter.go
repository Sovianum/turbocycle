package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type GasSplitter interface {
	graph.Node
	Input() nodes.ComplexGasSink
	MainOutput() nodes.ComplexGasSource
	ExtraOutput() nodes.ComplexGasSource
	ExtraWeight() float64
}

func NewGasSplitter(weight float64) GasSplitter {
	result := &gasSplitter{weight: weight}
	graph.AttachAllWithTags(
		result,
		[]*graph.Port{
			&result.gInput, &result.tInput, &result.pInput, &result.mrInput,
			&result.gOutput, &result.tOutput, &result.pOutput, &result.mrOutput,
			&result.gEOutput, &result.tEOutput, &result.pEOutput, &result.mrEOutput,
		},
		[]string{
			"gInput", "tInput", "pInput", "mrInput",
			"gOutput", "tOutput", "pOutput", "mrOutput",
			"gEOutput", "tEOutput", "pEOutput", "mrEOutput",
		},
	)
	return result
}

type gasSplitter struct {
	graph.BaseNode

	weight  float64
	tInput  graph.Port
	pInput  graph.Port
	gInput  graph.Port
	mrInput graph.Port

	tOutput  graph.Port
	pOutput  graph.Port
	gOutput  graph.Port
	mrOutput graph.Port

	tEOutput  graph.Port
	pEOutput  graph.Port
	gEOutput  graph.Port
	mrEOutput graph.Port
}

func (node *gasSplitter) GetName() string {
	return common.EitherString(node.GetInstanceName(), "GasSplitter")
}

func (node *gasSplitter) Process() error {
	imr := node.mrInput.GetState().Value().(float64)

	graph.SetAll(
		[]graph.PortState{
			node.gInput.GetState(),
			node.tInput.GetState(),
			node.pInput.GetState(),
			states.NewMassRatePortState(imr * (1 - node.weight)),
		},
		[]graph.Port{
			node.gOutput,
			node.tOutput,
			node.pOutput,
			node.mrOutput,
		},
	)

	graph.SetAll(
		[]graph.PortState{
			node.gInput.GetState(),
			node.tInput.GetState(),
			node.pInput.GetState(),
			states.NewMassRatePortState(imr * node.weight),
		},
		[]graph.Port{
			node.gEOutput,
			node.tEOutput,
			node.pEOutput,
			node.mrEOutput,
		},
	)
	return nil
}

func (node *gasSplitter) GetRequirePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.gInput, node.tInput, node.pInput, node.mrInput,
	}, nil
}

func (node *gasSplitter) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.gOutput, node.tOutput, node.pOutput, node.mrOutput,
		node.gEOutput, node.tEOutput, node.pEOutput, node.mrEOutput,
	}, nil
}

func (node *gasSplitter) GetPorts() []graph.Port {
	return []graph.Port{
		node.gInput, node.tInput, node.pInput, node.mrInput,
		node.gOutput, node.tOutput, node.pOutput, node.mrOutput,
		node.gEOutput, node.tEOutput, node.pEOutput, node.mrEOutput,
	}
}

func (node *gasSplitter) Input() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.gInput, node.tInput, node.pInput, node.mrInput,
	)
}

func (node *gasSplitter) MainOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.gOutput, node.tOutput, node.pOutput, node.mrOutput,
	)
}

func (node *gasSplitter) ExtraOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.gEOutput, node.tEOutput, node.pEOutput, node.mrEOutput,
	)
}

func (node *gasSplitter) ExtraWeight() float64 {
	return node.weight
}
