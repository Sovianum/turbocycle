package helper

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

type HubNode interface {
	graph.Node
	Inlet() graph.Port
	Outlet1() graph.Port
	Outlet2() graph.Port
}

func NewHubNode() HubNode {
	var result = &hubNode{}
	result.inlet = graph.NewAttachedPortWithTag(result, "inlet")
	result.outlet1 = graph.NewAttachedPortWithTag(result, "outlet1")
	result.outlet2 = graph.NewAttachedPortWithTag(result, "outlet2")
	return result
}

type hubNode struct {
	graph.BaseNode

	inlet   graph.Port
	outlet1 graph.Port
	outlet2 graph.Port
}

func (node *hubNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Hub")
}

func (node *hubNode) GetRequirePorts() ([]graph.Port, error) {
	return []graph.Port{node.inlet}, nil
}

func (node *hubNode) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{node.outlet1, node.outlet2}, nil
}

func (node *hubNode) GetPorts() []graph.Port {
	return []graph.Port{node.inlet, node.outlet1, node.outlet2}
}

func (node *hubNode) Process() error {
	node.outlet1.SetState(node.inlet.GetState())
	node.outlet2.SetState(node.inlet.GetState())
	return nil
}

func (node *hubNode) Inlet() graph.Port {
	return node.inlet
}

func (node *hubNode) Outlet1() graph.Port {
	return node.outlet1
}

func (node *hubNode) Outlet2() graph.Port {
	return node.outlet2
}
