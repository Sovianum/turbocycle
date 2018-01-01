package helper

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
)

type HubNode interface {
	core.Node
	Inlet() core.Port
	Outlet1() core.Port
	Outlet2() core.Port
}

func NewHubNode() HubNode {
	var result = &hubNode{}
	result.inlet = core.NewAttachedPort(result)
	result.outlet1 = core.NewAttachedPort(result)
	result.outlet2 = core.NewAttachedPort(result)
	return result
}

type hubNode struct {
	core.BaseNode

	inlet   core.Port
	outlet1 core.Port
	outlet2 core.Port
}

func (node *hubNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Hub")
}

func (node *hubNode) GetRequirePorts() []core.Port {
	return []core.Port{node.inlet}
}

func (node *hubNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.outlet1, node.outlet2}
}

func (node *hubNode) GetPorts() []core.Port {
	return []core.Port{node.inlet, node.outlet1, node.outlet2}
}

func (node *hubNode) Process() error {
	node.outlet1.SetState(node.inlet.GetState())
	node.outlet2.SetState(node.inlet.GetState())
	return nil
}

func (node *hubNode) Inlet() core.Port {
	return node.inlet
}

func (node *hubNode) Outlet1() core.Port {
	return node.outlet1
}

func (node *hubNode) Outlet2() core.Port {
	return node.outlet2
}
