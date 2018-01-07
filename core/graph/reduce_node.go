package graph

import "github.com/Sovianum/turbocycle/common"

type PortReduceFunc func(curr float64, port Port) (float64, error)
type FloatReduceFunc func(curr, new float64) (float64, error)

type ReduceNode interface {
	Node
	AddPortGroup(ports ...Port)
	OutputPort() Port
}

func NewReduceNode(groupReduceFunc PortReduceFunc, totalReduceFunc FloatReduceFunc, groupInit, totalInit float64) ReduceNode {
	var result = &reduceNode{
		inputPorts:  make([]Port, 0),
		pairPorts:   make([]Port, 0),
		groupLimits: []int{0},

		groupReduceFunc: groupReduceFunc,
		totalReduceFunc: totalReduceFunc,

		groupInit: groupInit,
		totalInit: totalInit,
	}
	result.outputPort = NewAttachedPort(result)
	return result
}

type reduceNode struct {
	BaseNode

	inputPorts  []Port
	pairPorts   []Port
	groupLimits []int

	outputPort Port

	groupReduceFunc PortReduceFunc
	totalReduceFunc FloatReduceFunc

	groupInit float64
	totalInit float64
}

func (node *reduceNode) OutputPort() Port {
	return node.outputPort
}

func (node *reduceNode) AddPortGroup(ports ...Port) {
	node.groupLimits = append(node.groupLimits, len(node.inputPorts)+len(ports))
	node.inputPorts = append(node.inputPorts, ports...)

	for _, port := range ports {
		var pairPort = NewAttachedPort(node)
		Link(port, pairPort)
		node.pairPorts = append(node.pairPorts, pairPort)
	}
}

func (node *reduceNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ZeroSumEquation")
}

func (node *reduceNode) Process() error {
	var limits = len(node.groupLimits)

	var totalResult = node.totalInit
	for i, j := 0, 1; i != limits-1 && j != limits; i, j = i+1, j+1 {
		var group = node.inputPorts[node.groupLimits[i]:node.groupLimits[j]]
		var groupResult = node.groupInit

		var err error = nil
		for _, port := range group {
			groupResult, err = node.groupReduceFunc(groupResult, port)
			if err != nil {
				return err
			}
		}

		totalResult, err = node.totalReduceFunc(totalResult, groupResult)
		if err != nil {
			return err
		}
	}

	node.outputPort.SetState(NewNumberPortState(totalResult))
	return nil
}

func (node *reduceNode) GetRequirePorts() []Port {
	return node.pairPorts
}

func (node *reduceNode) GetUpdatePorts() []Port {
	return []Port{node.outputPort}
}

func (node *reduceNode) GetPorts() []Port {
	return append(node.inputPorts, node.outputPort)
}
