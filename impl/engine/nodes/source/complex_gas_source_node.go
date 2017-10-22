package source

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type ComplexGasSourceNode interface {
	core.Node
	nodes.ComplexGasSource
}

type complexGasSourceNode struct {
	ports core.PortsType
	pStag float64
	tStag float64
	gas   gases.Gas
}

func NewComplexGasSourceNode(gas gases.Gas, tStag, pStag float64) ComplexGasSourceNode {
	var result = &complexGasSourceNode{
		ports: make(core.PortsType),
		pStag: pStag,
		tStag: tStag,
		gas:   gas,
	}

	result.ports[nodes.ComplexGasOutput] = core.NewPort()
	result.ports[nodes.ComplexGasOutput].SetInnerNode(result)
	result.ports[nodes.ComplexGasOutput].SetState(states.NewComplexGasPortState(gas, tStag, pStag, 1))

	return result
}

func (node *complexGasSourceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasOutputState core.PortState `json:"gas_output_state"`
	}{
		GasOutputState: node.ComplexGasOutput().GetState(),
	})
}

func (node *complexGasSourceNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *complexGasSourceNode) Process() error {
	node.ports[nodes.ComplexGasOutput].SetState(states.NewComplexGasPortState(node.gas, node.tStag, node.pStag, 1))
	return nil
}

func (node *complexGasSourceNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *complexGasSourceNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.ComplexGasOutput}, nil
}

func (node *complexGasSourceNode) GetPortTags() []string {
	return []string{nodes.ComplexGasOutput}
}

func (node *complexGasSourceNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasOutput:
		return node.ports[nodes.ComplexGasOutput], nil
	default:
		return nil, fmt.Errorf("Port %s of complexGasSourceNode can not be found", tag)
	}
}

func (node *complexGasSourceNode) ContextDefined() bool {
	return true
}

func (node *complexGasSourceNode) ComplexGasOutput() core.Port {
	return node.ports[nodes.ComplexGasOutput]
}

func (node *complexGasSourceNode) TStagOut() float64 {
	return node.ports[nodes.ComplexGasOutput].GetState().(states.ComplexGasPortState).TStag
}

func (node *complexGasSourceNode) PStagOut() float64 {
	return node.ports[nodes.ComplexGasOutput].GetState().(states.ComplexGasPortState).PStag
}
