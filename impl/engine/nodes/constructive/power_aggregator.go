package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

// PowerAggregator passes such labour to power output, that multiplied by total massRate
// of both its inputs is equal to total power (madness, i know)
// labour is normalized by main mass rate
// TODO pass power along power channels instead of labour
type PowerAggregator interface {
	graph.Node
	MainPowerInput() graph.Port
	MainMassRateInput() graph.Port
	ExtraPowerInput() graph.Port
	ExtraMassRateInput() graph.Port
	PowerOutput() graph.Port
}

func NewPowerAggregator() PowerAggregator {
	result := &powerAggregator{}
	graph.AttachAllWithTags(
		result,
		[]*graph.Port{
			&result.mainPowerInput, &result.mainMassRateInput,
			&result.extraPowerInput, &result.extraMassRateInput,
			&result.powerOutput,
		},
		[]string{
			"mainPowerInput", "mainMassRateInput",
			"extraPowerInput", "extraMassRateInput",
			"powerOutput",
		},
	)
	return result
}

type powerAggregator struct {
	graph.BaseNode

	mainPowerInput     graph.Port
	mainMassRateInput  graph.Port
	extraPowerInput    graph.Port
	extraMassRateInput graph.Port
	powerOutput        graph.Port
}

func (node *powerAggregator) MainPowerInput() graph.Port {
	return node.mainPowerInput
}

func (node *powerAggregator) MainMassRateInput() graph.Port {
	return node.mainMassRateInput
}

func (node *powerAggregator) ExtraPowerInput() graph.Port {
	return node.extraPowerInput
}

func (node *powerAggregator) ExtraMassRateInput() graph.Port {
	return node.extraMassRateInput
}

func (node *powerAggregator) PowerOutput() graph.Port {
	return node.powerOutput
}

func (node *powerAggregator) GetName() string {
	return common.EitherString(node.GetInstanceName(), "PowerAggregator")
}

func (node *powerAggregator) Process() error {
	mainMR := node.mainMassRateInput.GetState().Value().(float64)
	mainLabour := node.mainPowerInput.GetState().Value().(float64)

	extraMR := node.extraMassRateInput.GetState().Value().(float64)
	extraLabour := node.extraPowerInput.GetState().Value().(float64)

	node.powerOutput.SetState(states.NewPowerPortState(
		// divided by mainMR not by error, but to fulfill
		// energy balance according to current blocked stage
		// implementation
		(mainMR*mainLabour + extraMR*extraLabour) / mainMR,
	))
	return nil
}

func (node *powerAggregator) GetRequirePorts() ([]graph.Port, error) {
	return []graph.Port{
		node.mainPowerInput, node.mainMassRateInput,
		node.extraPowerInput, node.extraMassRateInput,
	}, nil
}

func (node *powerAggregator) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{node.powerOutput}, nil
}

func (node *powerAggregator) GetPorts() []graph.Port {
	return []graph.Port{
		node.mainPowerInput, node.mainMassRateInput,
		node.extraPowerInput, node.extraMassRateInput,
		node.powerOutput,
	}
}
