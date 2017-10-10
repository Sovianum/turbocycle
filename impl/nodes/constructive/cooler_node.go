package constructive

import (
	"encoding/json"
	"fmt"

	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
)

func NewCoolerNode(tOut float64, sigma float64) CoolerNode {
	var result = &coolerNode{
		ports: make(core.PortsType),
		tOut:  tOut,
		sigma: sigma,
	}

	result.ports[nodes.ComplexGasInput] = core.NewPort()
	result.ports[nodes.ComplexGasInput].SetInnerNode(result)
	result.ports[nodes.ComplexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.ComplexGasOutput] = core.NewPort()
	result.ports[nodes.ComplexGasOutput].SetInnerNode(result)
	result.ports[nodes.ComplexGasOutput].SetState(states.StandardAtmosphereState())

	return result
}

type CoolerNode interface {
	core.Node
	nodes.ComplexGasChannel
}

type coolerNode struct {
	ports core.PortsType
	tOut  float64
	sigma float64
}

func (node *coolerNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState  core.PortState `json:"gas_input_state"`
		GasOutputState core.PortState `json:"gas_output_state"`
		TOut           float64        `json:"t_out"`
	}{
		GasInputState:  node.complexGasInput().GetState(),
		GasOutputState: node.complexGasOutput().GetState(),
		TOut:           node.tOut,
	})
}

func (node *coolerNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *coolerNode) Process() error {
	var state = node.complexGasInput().GetState().(states.ComplexGasPortState)
	state.TStag = node.tOut
	state.PStag *= node.sigma
	node.complexGasOutput().SetState(state)
	return nil
}

func (node *coolerNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.ComplexGasInput}, nil
}

func (node *coolerNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.ComplexGasOutput}, nil
}

func (node *coolerNode) GetPortTags() []string {
	return []string{nodes.ComplexGasInput, nodes.ComplexGasOutput}
}

func (node *coolerNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasInput:
		return node.complexGasInput(), nil
	case nodes.ComplexGasOutput:
		return node.complexGasOutput(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found", tag)
	}
}

func (node *coolerNode) ContextDefined() bool {
	return true
}

func (node *coolerNode) ComplexGasInput() core.Port {
	return node.complexGasInput()
}

func (node *coolerNode) ComplexGasOutput() core.Port {
	return node.complexGasOutput()
}

func (node *coolerNode) complexGasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *coolerNode) complexGasOutput() core.Port {
	return node.ports[nodes.ComplexGasOutput]
}
