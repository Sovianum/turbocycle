package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
)

type GasSinkNode interface {
	core.Node
	ComplexGasSink
}

type gasSinkNode struct {
	ports core.PortsType
}

func NewGasSinkNode() GasSinkNode {
	var result = &gasSinkNode{
		ports: make(core.PortsType),
	}

	result.ports[complexGasInput] = core.NewPort()
	result.ports[complexGasInput].SetInnerNode(result)
	result.ports[complexGasInput].SetState(states.StandardAtmosphereState())

	return result
}

func (node *gasSinkNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState core.PortState `json:"gas_input_state"`
	}{
		GasInputState: node.ComplexGasInput().GetState(),
	})
}

func (node *gasSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *gasSinkNode) Process() error {
	return nil
}

func (node *gasSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{complexGasInput}, nil
}

func (node *gasSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *gasSinkNode) GetPortTags() []string {
	return []string{complexGasInput}
}

func (node *gasSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case complexGasInput:
		return node.ports[complexGasInput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of gasSinkNode can not be found", tag))
	}
}

func (node *gasSinkNode) ContextDefined() bool {
	return true
}

func (node *gasSinkNode) ComplexGasInput() core.Port {
	return node.ports[complexGasInput]
}

func (node *gasSinkNode) TStagIn() float64 {
	return node.ports[complexGasInput].GetState().(states.ComplexGasPortState).TStag
}

func (node *gasSinkNode) PStagIn() float64 {
	return node.ports[complexGasInput].GetState().(states.ComplexGasPortState).PStag
}
