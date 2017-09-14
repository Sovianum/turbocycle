package nodes

import (
	"github.com/Sovianum/turbocycle/core"
	"errors"
	"fmt"
)

type PowerSinkNode interface {
	core.Node
	PowerSink
}

type powerSinkNode struct {
	ports core.PortsType
}

func NewPortSinkNode() PowerSinkNode {
	var result = &powerSinkNode{
		ports:make(core.PortsType),
	}

	result.ports[powerInput] = core.NewPort()
	result.ports[powerInput].SetInnerNode(result)
	return result
}

func (node *powerSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *powerSinkNode) Process() error {
	return nil
}

func (node *powerSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *powerSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *powerSinkNode) GetPortTags() []string {
	return []string{powerInput}
}

func (node *powerSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case powerInput:
		return node.ports[gasInput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of powerSinkNode can not be found", tag))
	}
}

func (node *powerSinkNode) ContextDefined() bool {
	return true
}

func (node *powerSinkNode) PowerInput() core.Port {
	return node.ports[powerInput]
}



