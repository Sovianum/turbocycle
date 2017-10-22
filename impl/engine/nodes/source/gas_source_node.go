package source

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type GasSourceNode interface {
	core.Node
	nodes.GasSource
}

type gasSourceNode struct {
	ports core.PortsType
	gas   gases.Gas
}

func NewGasSourceNode(gas gases.Gas) GasSourceNode {
	var result = &gasSourceNode{
		ports: make(core.PortsType),
		gas:   gas,
	}

	result.ports[nodes.GasOutput] = core.NewPort()
	result.ports[nodes.GasOutput].SetInnerNode(result)
	result.ports[nodes.GasOutput].SetState(states.NewGasPortState(gas))

	return result
}

func (node *gasSourceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasOutputState core.PortState `json:"gas_output_state"`
	}{
		GasOutputState: node.GasOutput().GetState(),
	})
}

func (node *gasSourceNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *gasSourceNode) Process() error {
	node.ports[nodes.GasOutput].SetState(states.NewGasPortState(node.gas))
	return nil
}

func (node *gasSourceNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *gasSourceNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.GasOutput}, nil
}

func (node *gasSourceNode) GetPortTags() []string {
	return []string{nodes.GasOutput}
}

func (node *gasSourceNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.GasOutput:
		return node.ports[nodes.GasOutput], nil
	default:
		return nil, fmt.Errorf("Port %s of gasSourceNode can not be found", tag)
	}
}

func (node *gasSourceNode) ContextDefined() bool {
	return true
}

func (node *gasSourceNode) GasOutput() core.Port {
	return node.ports[nodes.GasOutput]
}
