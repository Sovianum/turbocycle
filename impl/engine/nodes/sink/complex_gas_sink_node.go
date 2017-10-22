package sink

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type ComplexGasSinkNode interface {
	core.Node
	nodes.ComplexGasSink
}

type complexGasSinkNode struct {
	ports core.PortsType
}

func NewComplexGasSinkNode() ComplexGasSinkNode {
	var result = &complexGasSinkNode{
		ports: make(core.PortsType),
	}

	result.ports[nodes.ComplexGasInput] = core.NewPort()
	result.ports[nodes.ComplexGasInput].SetInnerNode(result)
	result.ports[nodes.ComplexGasInput].SetState(states.StandardAtmosphereState())

	return result
}

func (node *complexGasSinkNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState core.PortState `json:"gas_input_state"`
	}{
		GasInputState: node.ComplexGasInput().GetState(),
	})
}

func (node *complexGasSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *complexGasSinkNode) Process() error {
	return nil
}

func (node *complexGasSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.ComplexGasInput}, nil
}

func (node *complexGasSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *complexGasSinkNode) GetPortTags() []string {
	return []string{nodes.ComplexGasInput}
}

func (node *complexGasSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasInput:
		return node.ports[nodes.ComplexGasInput], nil
	default:
		return nil, fmt.Errorf("Port %s of complexGasSinkNode can not be found", tag)
	}
}

func (node *complexGasSinkNode) ContextDefined() bool {
	return true
}

func (node *complexGasSinkNode) ComplexGasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *complexGasSinkNode) TStagIn() float64 {
	return node.ports[nodes.ComplexGasInput].GetState().(states.ComplexGasPortState).TStag
}

func (node *complexGasSinkNode) PStagIn() float64 {
	return node.ports[nodes.ComplexGasInput].GetState().(states.ComplexGasPortState).PStag
}
