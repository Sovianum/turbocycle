package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type TransmissionNode interface {
	core.Node
	nodes.PowerChannel
	Eta() float64
}

func NewTransmissionNode(etaM float64) TransmissionNode {
	var result = &transmissionNode{
		etaM: etaM,
	}

	result.powerInput = core.NewAttachedPort(result)
	result.powerOutput = core.NewAttachedPort(result)

	return result
}

type transmissionNode struct {
	core.BaseNode

	powerInput  core.Port
	powerOutput core.Port

	etaM float64
}

func (node *transmissionNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Transmission")
}

func (node *transmissionNode) GetPorts() []core.Port {
	return []core.Port{node.powerInput, node.powerOutput}
}

func (node *transmissionNode) GetRequirePorts() []core.Port {
	return []core.Port{node.powerInput}
}

func (node *transmissionNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.powerOutput}
}

func (node *transmissionNode) Eta() float64 {
	return node.etaM
}

func (node *transmissionNode) PowerInput() core.Port {
	return node.powerInput
}

func (node *transmissionNode) PowerOutput() core.Port {
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
