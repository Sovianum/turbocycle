package nodes

import (
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
)

type PowerSourceNode interface {
	core.Node
	PowerSource
}

type powerSourceNode struct {
	ports     core.PortsType
	lSpecific float64
}

func NewPowerSourceNode(lRel float64) PowerSourceNode {
	var result = &powerSourceNode{
		ports:     make(core.PortsType),
		lSpecific: lRel,
	}

	result.ports[powerOutput] = core.NewPort()
	result.ports[powerOutput].SetInnerNode(result)
	result.ports[powerOutput].SetState(states.StandartPowerState())

	return result
}

func (node *powerSourceNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *powerSourceNode) Process() error {
	node.ports[powerOutput].SetState(states.NewPowerPortState(node.lSpecific))
	return nil
}

func (node *powerSourceNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *powerSourceNode) GetUpdatePortTags() ([]string, error) {
	return []string{powerOutput}, nil
}

func (node *powerSourceNode) GetPortTags() []string {
	return []string{powerOutput}
}

func (node *powerSourceNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case powerOutput:
		return node.ports[powerOutput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of powerSourceNode can not be found", tag))
	}
}

func (node *powerSourceNode) ContextDefined() bool {
	return true
}

func (node *powerSourceNode) PowerOutput() core.Port {
	return node.ports[powerOutput]
}
