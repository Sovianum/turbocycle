package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
)

type GasSourceNode interface {
	core.Node
	ComplexGasSource
}

type gasSourceNode struct {
	ports core.PortsType
	pStag float64
	tStag float64
	gas   gases.Gas
}

func NewGasSource(gas gases.Gas, tStag, pStag float64) GasSourceNode {
	var result = &gasSourceNode{
		ports: make(core.PortsType),
		pStag: pStag,
		tStag: tStag,
		gas:   gas,
	}

	result.ports[complexGasOutput] = core.NewPort()
	result.ports[complexGasOutput].SetInnerNode(result)
	result.ports[complexGasOutput].SetState(states.StandardAtmosphereState())

	return result
}

func (node *gasSourceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasOutputState core.PortState `json:"gas_output_state"`
	}{
		GasOutputState: node.ComplexGasOutput().GetState(),
	})
}

func (node *gasSourceNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *gasSourceNode) Process() error {
	node.ports[complexGasOutput].SetState(states.NewComplexGasPortState(node.gas, node.tStag, node.pStag, 1))
	return nil
}

func (node *gasSourceNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *gasSourceNode) GetUpdatePortTags() ([]string, error) {
	return []string{complexGasOutput}, nil
}

func (node *gasSourceNode) GetPortTags() []string {
	return []string{complexGasOutput}
}

func (node *gasSourceNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case complexGasOutput:
		return node.ports[complexGasOutput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of gasSourceNode can not be found", tag))
	}
}

func (node *gasSourceNode) ContextDefined() bool {
	return true
}

func (node *gasSourceNode) ComplexGasOutput() core.Port {
	return node.ports[complexGasOutput]
}

func (node *gasSourceNode) TStagOut() float64 {
	return node.ports[complexGasOutput].GetState().(states.ComplexGasPortState).TStag
}

func (node *gasSourceNode) PStagOut() float64 {
	return node.ports[complexGasInput].GetState().(states.ComplexGasPortState).PStag
}
