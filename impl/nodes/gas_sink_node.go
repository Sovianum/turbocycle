package nodes

import (
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
)

type GasSinkNode interface {
	core.Node
	GasSink
}

type gasSinkNode struct {
	ports core.PortsType
}

func NewGasSinkNode() GasSinkNode {
	var result = &gasSinkNode{
		ports: make(core.PortsType),
	}

	result.ports[gasInput] = core.NewPort()
	result.ports[gasInput].SetInnerNode(result)
	result.ports[gasInput].SetState(states.StandartAtmosphereState())

	return result
}

func (node *gasSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *gasSinkNode) Process() error {
	return nil
}

func (node *gasSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{gasInput}, nil
}

func (node *gasSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *gasSinkNode) GetPortTags() []string {
	return []string{gasInput}
}

func (node *gasSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case gasInput:
		return node.ports[gasInput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of gasSinkNode can not be found", tag))
	}
}

func (node *gasSinkNode) ContextDefined() bool {
	return true
}

func (node *gasSinkNode) GasInput() core.Port {
	return node.ports[gasInput]
}

func (node *gasSinkNode) TStagIn() float64 {
	return node.ports[gasInput].GetState().(states.GasPortState).TStag
}

func (node *gasSinkNode) PStagIn() float64 {
	return node.ports[gasInput].GetState().(states.GasPortState).PStag
}
