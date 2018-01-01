package source

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type PressureSourceNode interface {
	graph.Node
	nodes.PressureSource
}

func NewPressureSourceNode(pStag float64) PressureSourceNode {
	var result = &pressureSourceNode{
		pStag: pStag,
	}

	result.setOutput(graph.NewAttachedPort(result))
	return result
}

type pressureSourceNode struct {
	sourceNode
	pStag float64
}

func (node *pressureSourceNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "PressureSource")
}

func (node *pressureSourceNode) Process() error {
	node.getOutput().SetState(states.NewPressurePortState(node.pStag))
	return nil
}

func (node *pressureSourceNode) PStagOut() float64 {
	return node.pStag
}

func (node *pressureSourceNode) PressureOutput() graph.Port {
	return node.getOutput()
}
