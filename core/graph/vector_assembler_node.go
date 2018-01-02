package graph

import (
	"sort"

	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
)

// this node collects values from its inputs and assembles it to
// vector. the order of adding is preserved
type VectorAssemblerNode interface {
	Node
	GetPairPort(outerPort Port) Port
	DeletePairPort(port Port)
	GetVectorPort() Port
}

func NewVectorAssemblerNode() VectorAssemblerNode {
	return &vectorAssemblerNode{
		portCnt:      0,
		inputPortMap: make(map[int]Port),
		outerPortMap: make(map[Port]int),
		vectorPort:   NewPort(),
	}
}

type vectorAssemblerNode struct {
	BaseNode

	portCnt      int
	inputPortMap map[int]Port
	outerPortMap map[Port]int
	vectorPort   Port
}

func (node *vectorAssemblerNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "VectorAssembler")
}

func (node *vectorAssemblerNode) Process() error {
	var inputPorts = node.getInputPorts()
	var values = make([]float64, len(inputPorts))

	for i, port := range inputPorts {
		values[i] = port.GetState().(NumberPortState).Num
	}

	node.vectorPort.SetState(NewVectorPortState(
		mat.NewVecDense(len(values), values),
	))
	return nil
}

func (node *vectorAssemblerNode) GetRequirePorts() []Port {
	return node.getInputPorts()
}

func (node *vectorAssemblerNode) GetUpdatePorts() []Port {
	return []Port{node.vectorPort}
}

func (node *vectorAssemblerNode) GetPorts() []Port {
	var ports = node.getInputPorts()
	ports = append(ports, node.vectorPort)
	return ports
}

func (node *vectorAssemblerNode) GetPairPort(outerPort Port) Port {
	var portId, ok = node.outerPortMap[outerPort]
	if ok {
		return node.inputPortMap[portId]
	}

	node.portCnt++
	node.outerPortMap[outerPort] = node.portCnt
	var port = NewPort()
	node.inputPortMap[node.portCnt] = port

	return port
}

func (node *vectorAssemblerNode) DeletePairPort(port Port) {
	var id, ok = node.outerPortMap[port]
	if !ok {
		return
	}
	delete(node.outerPortMap, port)
	delete(node.inputPortMap, id)
}

func (node *vectorAssemblerNode) GetVectorPort() Port {
	return node.vectorPort
}

func (node *vectorAssemblerNode) getInputPorts() []Port {
	var keys = make([]int, 0)
	for key := range node.inputPortMap {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	var ports = make([]Port, len(keys))
	for i, key := range keys {
		ports[i] = node.inputPortMap[key]
	}
	return ports
}
