package helper

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
)

type CycleBreakNode interface {
	core.Node
	UpdatePort() core.Port
	DataSourcePort() core.Port
}

func NewCycleBreakerNode() CycleBreakNode {
	var result = &cycleBreakNode{}
	result.updatePort = core.NewAttachedPort(result)
	result.sourcePort = core.NewAttachedPort(result)
	return result
}

type cycleBreakNode struct {
	core.BaseNode
	updatePort core.Port
	sourcePort core.Port
}

func (node *cycleBreakNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "CycleBreak")
}

func (node *cycleBreakNode) GetRequirePorts() []core.Port {
	return make([]core.Port, 0)
}

func (node *cycleBreakNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.updatePort}
}

func (node *cycleBreakNode) GetPorts() []core.Port {
	return []core.Port{node.sourcePort, node.updatePort}
}

func (node *cycleBreakNode) Process() error {
	node.updatePort.SetState(node.sourcePort.GetState())
	return nil
}

func (node *cycleBreakNode) UpdatePort() core.Port {
	return node.updatePort
}

func (node *cycleBreakNode) DataSourcePort() core.Port {
	return node.sourcePort
}
