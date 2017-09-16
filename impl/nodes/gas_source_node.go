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
	GasSource
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

	result.ports[gasOutput] = core.NewPort()
	result.ports[gasOutput].SetInnerNode(result)
	result.ports[gasOutput].SetState(states.StandardAtmosphereState())

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
	node.ports[gasOutput].SetState(states.NewGasPortState(node.gas, node.tStag, node.pStag, 1))
	return nil
}

func (node *gasSourceNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *gasSourceNode) GetUpdatePortTags() ([]string, error) {
	return []string{gasOutput}, nil
}

func (node *gasSourceNode) GetPortTags() []string {
	return []string{gasOutput}
}

func (node *gasSourceNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case gasOutput:
		return node.ports[gasOutput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of gasSourceNode can not be found", tag))
	}
}

func (node *gasSourceNode) ContextDefined() bool {
	return true
}

func (node *gasSourceNode) GasOutput() core.Port {
	return node.ports[gasOutput]
}

func (node *gasSourceNode) TStagOut() float64 {
	return node.ports[gasOutput].GetState().(states.GasPortState).TStag
}

func (node *gasSourceNode) PStagOut() float64 {
	return node.ports[gasInput].GetState().(states.GasPortState).PStag
}
