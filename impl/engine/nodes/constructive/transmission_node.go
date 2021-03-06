package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

// Transmission turbine transmits energy DEMAND, not the energy supply
// so its power input must be attached to compressor
type TransmissionNode interface {
	graph.Node
	nodes.PowerChannel
	Eta() float64
}

func NewTransmissionNode(etaM float64) TransmissionNode {
	var result = &transmissionNode{
		etaM: etaM,
	}

	result.powerInput = graph.NewAttachedPortWithTag(result, nodes.PowerInputTag)
	result.powerOutput = graph.NewAttachedPortWithTag(result, nodes.PowerOutputTag)

	return result
}

type transmissionNode struct {
	graph.BaseNode

	powerInput  graph.Port
	powerOutput graph.Port

	etaM float64
}

func (node *transmissionNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Transmission")
}

func (node *transmissionNode) GetPorts() []graph.Port {
	return []graph.Port{node.powerInput, node.powerOutput}
}

func (node *transmissionNode) GetRequirePorts() ([]graph.Port, error) {
	return []graph.Port{node.powerInput}, nil
}

func (node *transmissionNode) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{node.powerOutput}, nil
}

func (node *transmissionNode) Eta() float64 {
	return node.etaM
}

func (node *transmissionNode) PowerInput() graph.Port {
	return node.powerInput
}

func (node *transmissionNode) PowerOutput() graph.Port {
	return node.powerOutput
}

func (node *transmissionNode) Process() error {
	var inputState = node.PowerInput().GetState()

	switch v := inputState.(type) {
	case states.PowerPortState:
		var state = inputState.(states.PowerPortState)
		state.LSpecific /= node.etaM // TODO check if division is a correct operation (depends on ifo flow direction)
		node.powerOutput.SetState(state)
		return nil
	default:
		return common.GetTypeError("PowerPortState", v)
	}
}
