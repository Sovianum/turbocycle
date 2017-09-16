package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/Sovianum/turbocycle/impl/nodes"
)

type PressureSourceNode interface {
	core.Node
	nodes.PressureSource
}

type pressureSourceNode struct {
	ports core.PortsType
}

func NewPressureSinkNode() PressureSourceNode {
	var result = &pressureSourceNode{
		ports: make(core.PortsType),
	}

	result.ports[nodes.PressureOutput] = core.NewPort()
	result.ports[nodes.PressureOutput].SetInnerNode(result)
	result.ports[nodes.PressureOutput].SetState(states.NewPressurePortState(1e5))

	return result
}

func (node *pressureSourceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PressureOutputState core.PortState `json:"pressure_output_state"`
	}{
		PressureOutputState: node.PressureOutput().GetState(),
	})
}

func (node *pressureSourceNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *pressureSourceNode) Process() error {
	return nil
}

func (node *pressureSourceNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.PressureOutput}, nil
}

func (node *pressureSourceNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *pressureSourceNode) GetPortTags() []string {
	return []string{nodes.PressureOutput}
}

func (node *pressureSourceNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.PressureOutput:
		return node.ports[nodes.PressureOutput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of pressureSourceNode can not be found", tag))
	}
}

func (node *pressureSourceNode) ContextDefined() bool {
	return true
}

func (node *pressureSourceNode) PStagOut() float64 {
	return node.ports[nodes.PressureOutput].GetState().(states.PressurePortState).PStag
}

func (node *pressureSourceNode) PressureOutput() core.Port {
	return node.ports[nodes.PressureOutput]
}