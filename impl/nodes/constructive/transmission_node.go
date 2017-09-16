package constructive

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
)

type TransmissionNode interface {
	core.Node
	nodes.PowerChannel
}

type transmissionNode struct {
	ports core.PortsType
	etaM  float64
}

func NewTransmissionNode(etaM float64) TransmissionNode {
	var transmissionNode = &transmissionNode{
		ports: make(core.PortsType),
		etaM:  etaM,
	}

	var inputPort = core.NewPort()
	inputPort.SetInnerNode(transmissionNode)
	transmissionNode.ports[nodes.PowerInput] = inputPort
	transmissionNode.ports[nodes.PowerInput].SetState(states.StandardPowerState())

	var outputPort = core.NewPort()
	outputPort.SetInnerNode(transmissionNode)
	transmissionNode.ports[nodes.PowerOutput] = outputPort
	transmissionNode.ports[nodes.PowerOutput].SetState(states.StandardPowerState())

	return transmissionNode
}

func (node *transmissionNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PowerInputState  core.PortState `json:"power_input_state"`
		PowerOutputState core.PortState `json:"power_output_state"`
		EtaM             float64        `json:"eta_m"`
	}{
		PowerInputState:  node.powerInput().GetState(),
		PowerOutputState: node.powerOutput().GetState(),
		EtaM:             node.etaM,
	})
}

func (node *transmissionNode) ContextDefined() bool {
	return true
}

func (node *transmissionNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.PowerInput:
		return node.powerInput(), nil
	case nodes.PowerOutput:
		return node.PowerOutput(), nil
	default:
		return nil, errors.New(fmt.Sprintf("port with tag \"%s\" not found", tag))
	}
}

func (node *transmissionNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.PowerInput}, nil
}

func (node *transmissionNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.PowerOutput}, nil
}

func (node *transmissionNode) GetPortTags() []string {
	return []string{nodes.PowerInput, nodes.PowerOutput}
}

func (node *transmissionNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *transmissionNode) PowerInput() core.Port {
	return node.powerInput()
}

func (node *transmissionNode) PowerOutput() core.Port {
	return node.powerOutput()
}

func (node *transmissionNode) Process() error {
	var inputState = node.PowerInput().GetState()

	switch v := inputState.(type) {
	case states.PowerPortState:
		var state = inputState.(states.PowerPortState)
		state.LSpecific /= node.etaM // TODO check if division is a correct operation (depends on ifo flow direction)
		node.powerOutput().SetState(state)
		return nil
	default:
		return common.GetTypeError("PowerPortState", v)
	}
}

func (node *transmissionNode) powerInput() core.Port {
	return node.ports[nodes.PowerInput]
}

func (node *transmissionNode) powerOutput() core.Port {
	return node.ports[nodes.PowerOutput]
}
