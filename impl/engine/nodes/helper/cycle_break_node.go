package helper

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

type CycleBreakNode interface {
	graph.Node
	UpdatePort() graph.Port
	DataSourcePort() graph.Port
}

func NewCycleBreakNode(initialState graph.PortState) CycleBreakNode {
	var result = &cycleBreakNode{}
	result.updatePort = graph.NewAttachedPortWithTag(result, "updatePort")
	result.sourcePort = graph.NewAttachedPortWithTag(result, "sourcePort")
	result.sourcePort.SetState(initialState)
	return result
}

type cycleBreakNode struct {
	graph.BaseNode
	updatePort graph.Port
	sourcePort graph.Port
}

func (node *cycleBreakNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "CycleBreak")
}

func (node *cycleBreakNode) GetRequirePorts() ([]graph.Port, error) {
	return make([]graph.Port, 0), nil
}

func (node *cycleBreakNode) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{node.updatePort}, nil
}

func (node *cycleBreakNode) GetPorts() []graph.Port {
	return []graph.Port{node.sourcePort, node.updatePort}
}

func (node *cycleBreakNode) Process() error {
	node.updatePort.SetState(node.sourcePort.GetState())
	return nil
}

func (node *cycleBreakNode) UpdatePort() graph.Port {
	return node.updatePort
}

func (node *cycleBreakNode) DataSourcePort() graph.Port {
	return node.sourcePort
}
