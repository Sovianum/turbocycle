package helper

import (
	"encoding/json"
	"github.com/Sovianum/turbocycle/core"
	"errors"
	"fmt"
)

const (
	inlet = "inlet"
	outlet1 = "outlet1"
	outlet2 = "outlet2"
)

type HubNode interface {
	core.Node
	Inlet() core.Port
	Outlet1() core.Port
	Outlet2() core.Port
}

type hubNode struct {
	ports core.PortsType
}

func NewHubNode(initState core.PortState) HubNode {
	var result = &hubNode{
		ports:make(core.PortsType),
	}

	result.ports[inlet] = core.NewPort()
	result.ports[inlet].SetInnerNode(result)
	result.ports[inlet].SetState(initState)

	result.ports[outlet1] = core.NewPort()
	result.ports[outlet1].SetInnerNode(result)
	result.ports[outlet1].SetState(initState)

	result.ports[outlet2] = core.NewPort()
	result.ports[outlet2].SetInnerNode(result)
	result.ports[outlet2].SetState(initState)

	return result
}

func (node *hubNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		State core.PortState
	}{
		State: node.Inlet().GetState(),
	})
}

func (node *hubNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *hubNode) Process() error {
	node.ports[outlet1].SetState(node.ports[inlet].GetState())
	node.ports[outlet2].SetState(node.ports[inlet].GetState())
	return nil
}

func (node *hubNode) GetRequirePortTags() ([]string, error) {
	return []string{inlet}, nil
}

func (node *hubNode) GetUpdatePortTags() ([]string, error) {
	return []string{outlet1, outlet2}, nil
}

func (node *hubNode) GetPortTags() []string {
	return []string{inlet, outlet1, outlet2}
}

func (node *hubNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case inlet:
		return node.ports[inlet], nil
	case outlet1:
		return node.ports[outlet1], nil
	case outlet2:
		return node.ports[outlet2], nil
	default:
		return nil, errors.New(fmt.Sprintf("tag \"%s\" was not found in hubNode", tag))
	}
}

func (node *hubNode) ContextDefined() bool {
	return true
}

func (node *hubNode) Inlet() core.Port {
	return node.ports[inlet]
}

func (node *hubNode) Outlet1() core.Port {
	return node.ports[outlet1]
}

func (node *hubNode) Outlet2() core.Port {
	return node.ports[outlet2]
}
