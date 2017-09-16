package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
)

type PowerSinkNode interface {
	core.Node
	PowerSink
}

type powerSinkNode struct {
	ports core.PortsType
}

func NewPortSinkNode() PowerSinkNode {
	var result = &powerSinkNode{
		ports: make(core.PortsType),
	}

	result.ports[powerInput] = core.NewPort()
	result.ports[powerInput].SetInnerNode(result)
	result.ports[powerInput].SetState(states.StandardPowerState())

	return result
}

func (node *powerSinkNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PowerInputState core.PortState `json:"power_input_state"`
	}{
		PowerInputState: node.PowerInput().GetState(),
	})
}

func (node *powerSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *powerSinkNode) Process() error {
	return nil
}

func (node *powerSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *powerSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *powerSinkNode) GetPortTags() []string {
	return []string{powerInput}
}

func (node *powerSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case powerInput:
		return node.ports[powerInput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of powerSinkNode can not be found", tag))
	}
}

func (node *powerSinkNode) ContextDefined() bool {
	return true
}

func (node *powerSinkNode) PowerInput() core.Port {
	return node.ports[powerInput]
}
