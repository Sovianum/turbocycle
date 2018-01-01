package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type ComplexGasSinkNode interface {
	graph.Node
	nodes.ComplexGasSink
}

func NewComplexGasSinkNode() ComplexGasSinkNode {
	var result = &complexGasSinkNode{}
	result.setInput(graph.NewAttachedPort(result))
	return result
}

type complexGasSinkNode struct {
	sinkNode
}

func (node *complexGasSinkNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ComplexGasSink")
}

func (node *complexGasSinkNode) ComplexGasInput() graph.Port {
	return node.getInput()
}

func (node *complexGasSinkNode) TStagIn() float64 {
	return node.getInput().GetState().(states.ComplexGasPortState).TStag
}

func (node *complexGasSinkNode) PStagIn() float64 {
	return node.getInput().GetState().(states.ComplexGasPortState).PStag
}
