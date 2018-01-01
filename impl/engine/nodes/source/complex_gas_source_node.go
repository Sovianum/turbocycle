package source

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type ComplexGasSourceNode interface {
	graph.Node
	nodes.ComplexGasSource
}

func NewComplexGasSourceNode(gas gases.Gas, tStag, pStag float64) ComplexGasSourceNode {
	var result = &complexGasSourceNode{
		pStag: pStag,
		tStag: tStag,
		gas:   gas,
	}

	result.setOutput(graph.NewAttachedPort(result))
	return result
}

type complexGasSourceNode struct {
	sourceNode
	pStag float64
	tStag float64
	gas   gases.Gas
}

func (node *complexGasSourceNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ComplexGasSource")
}

func (node *complexGasSourceNode) Process() error {
	node.getOutput().SetState(states.NewComplexGasPortState(node.gas, node.tStag, node.pStag, 1))
	return nil
}

func (node *complexGasSourceNode) ComplexGasOutput() graph.Port {
	return node.getOutput()
}

func (node *complexGasSourceNode) TStagOut() float64 {
	return node.tStag
}

func (node *complexGasSourceNode) PStagOut() float64 {
	return node.pStag
}
