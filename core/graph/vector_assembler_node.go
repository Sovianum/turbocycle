package graph

import (
	"sort"

	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
)

type NamedValue struct {
	Name  string
	Value float64
}

// this node collects values from its inputs and assembles it to
// vector. the order of adding is preserved
type VectorAssemblerNode interface {
	Node
	AddInputPorts(outerPorts ...Port)
	DeleteInputPorts(outerPorts ...Port)
	GetVectorPort() Port
	GetNamedReport() []NamedValue
}

func NewVectorAssemblerNode() VectorAssemblerNode {
	var node = &vectorAssemblerNode{
		portCnt:      0,
		inputPortMap: make(map[int]Port),
		outerPortMap: make(map[Port]int),
	}
	node.vectorPort = NewAttachedPort(node)
	return node
}

type vectorAssemblerNode struct {
	BaseNode

	portCnt      int
	inputPortMap map[int]Port
	outerPortMap map[Port]int
	vectorPort   Port
}

func (node *vectorAssemblerNode) GetNamedReport() []NamedValue {
	result := make([]NamedValue, node.portCnt)
	cnt := 0
	for _, port := range node.inputPortMap {
		result[cnt] = NamedValue{
			Name:  port.GetOuterNode().GetName(),
			Value: port.GetState().Value().(float64),
		}
		cnt++
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
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

func (node *vectorAssemblerNode) GetRequirePorts() ([]Port, error) {
	return node.getInputPorts(), nil
}

func (node *vectorAssemblerNode) GetUpdatePorts() ([]Port, error) {
	return []Port{node.vectorPort}, nil
}

func (node *vectorAssemblerNode) GetPorts() []Port {
	var ports = node.getInputPorts()
	ports = append(ports, node.vectorPort)
	return ports
}

func (node *vectorAssemblerNode) AddInputPorts(outerPorts ...Port) {
	for _, port := range outerPorts {
		node.addInputPort(port)
	}
}

func (node *vectorAssemblerNode) DeleteInputPorts(outerPorts ...Port) {
	for _, port := range outerPorts {
		node.deleteInputPort(port)
	}
}

func (node *vectorAssemblerNode) GetVectorPort() Port {
	return node.vectorPort
}

func (node *vectorAssemblerNode) addInputPort(outerPort Port) {
	var _, ok = node.outerPortMap[outerPort]
	if ok {
		return
	}

	node.portCnt++
	node.outerPortMap[outerPort] = node.portCnt
	var port = NewAttachedPort(node)
	node.inputPortMap[node.portCnt] = port
	Link(port, outerPort)
}

func (node *vectorAssemblerNode) deleteInputPort(outerPort Port) {
	var id, ok = node.outerPortMap[outerPort]
	if !ok {
		return
	}
	delete(node.outerPortMap, outerPort)
	delete(node.inputPortMap, id)
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
