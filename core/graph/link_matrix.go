package graph

import (
	"fmt"
	"strings"

	"github.com/Sovianum/turbocycle/common"
)

const (
	inaccessibleNodesMsg  = "inaccessible nodes detected"
	unconnectedPortsMsg   = "unconnected port detected"
	contextUndefinedNodes = "context undefined nodes detected"
)

func newGraphMatrix(nodeArr []Node) (*graphMatrix, GraphError) {
	var nodeMap = common.NewBiMap()
	for i, node := range nodeArr {
		nodeMap.Add(i, node)
	}

	var matrix = make([][]bool, len(nodeArr))
	var matrixCopy = make([][]bool, len(nodeArr))

	for i := 0; i != len(nodeArr); i++ {
		matrix[i] = make([]bool, len(nodeArr))
		matrixCopy[i] = make([]bool, len(nodeArr))
	}

	var result = &graphMatrix{
		nodes:      nodeMap,
		matrix:     matrix,
		matrixCopy: matrixCopy,
	}

	var err = result.setEdges()
	if err != nil {
		return nil, err
	}

	result.copyMatrix()
	return result, nil
}

// this type is used to perform graph related operations like
// getting call order, checking context definition and unconnected ports
type graphMatrix struct {
	nodes      common.BiMap
	matrix     [][]bool
	matrixCopy [][]bool
}

func (m *graphMatrix) GetCallOrder() ([]Node, GraphError) {
	m.copyMatrix()
	var usedFreeNodes = make(map[Node]bool)

	var getUnusedFreeNodes = func() []Node {
		var freeNodes = m.getFreeNodes()
		var result = make([]Node, 0)

		for _, node := range freeNodes {
			var _, ok = usedFreeNodes[node]
			if !ok {
				usedFreeNodes[node] = true
				result = append(result, node)
			}
		}

		return result
	}

	var result = make([]Node, 0)
	var front = getUnusedFreeNodes()

	for len(front) != 0 {
		result = append(result, front...)
		m.removeAllDependencies(front)
		front = getUnusedFreeNodes()
	}

	var dependentNodes = m.getDependentNodes()
	if len(dependentNodes) > 0 {
		var nodeNames = make([]string, len(dependentNodes))
		for i, node := range dependentNodes {
			nodeNames[i] = node.GetName()
		}
		return nil, graphErrorFromNodes(getInaccessibleErrMsg(dependentNodes), dependentNodes)
	}
	return result, nil
}

func (m *graphMatrix) copyMatrix() {
	for i, row := range m.matrix {
		for j, col := range row {
			m.matrixCopy[i][j] = col
		}
	}
}

// edges are set in require direction, i.e. if A requires B
// there exists edge from B to A, and matrix[i_A][i_B] == true
func (m *graphMatrix) setEdges() GraphError {
	if unconnectedPorts := m.getUnconnectedPorts(); len(unconnectedPorts) > 0 {
		var errMsg = getUnconnectedErrMsg(unconnectedPorts)
		return graphErrorFromPorts(errMsg, unconnectedPorts)
	}
	if undefined := m.getContextUndefinedNodes(); len(undefined) > 0 {
		var errMsg = getUndefinedErrMsg(undefined)
		return graphErrorFromNodes(errMsg, undefined)
	}

	for pair := range m.nodes.Iterate() {
		var innerNode = pair.Val.(Node)
		var requirePorts, portErr = innerNode.GetRequirePorts()
		if portErr != nil {
			return graphErrorFromPorts(portErr.Error(), nil)
		}

		for _, port := range requirePorts {
			var outerNode = port.GetOuterNode()

			var innerNodeId, _ = m.nodes.GetByVal(innerNode)
			var outerNodeId, _ = m.nodes.GetByVal(outerNode)
			m.matrix[innerNodeId][outerNodeId] = true
		}
	}
	return nil
}

func (m *graphMatrix) getUnconnectedPorts() []Port {
	var unconnected = make([]Port, 0)
	for pair := range m.nodes.Iterate() {
		for _, port := range pair.Val.(Node).GetPorts() {
			if port.GetOuterNode() == nil {
				unconnected = append(unconnected, port)
			}
		}
	}
	return unconnected
}

func (m *graphMatrix) getContextUndefinedNodes() []Node {
	var undefined = make([]Node, 0)
	for pair := range m.nodes.Iterate() {
		if node := pair.Val.(Node); !node.ContextDefined() {
			undefined = append(undefined, node)
		}
	}
	return undefined
}

// get nodes with at least one dependency
func (m *graphMatrix) getDependentNodes() []Node {
	var nodes = make(map[Node]bool)
	// add all nodes to the map
	for pair := range m.nodes.Iterate() {
		nodes[pair.Val.(Node)] = true
	}

	var freeNodes = m.getFreeNodes()
	// remove free nodes from the map
	for _, node := range freeNodes {
		delete(nodes, node)
	}

	var result = make([]Node, 0)
	for node := range nodes {
		result = append(result, node)
	}
	return result
}

func (m *graphMatrix) removeAllDependencies(nodes []Node) {
	for _, node := range nodes {
		m.removeDependency(node)
	}
}

func (m *graphMatrix) removeDependency(node Node) {
	var key, ok = m.nodes.GetByVal(node)
	if !ok {
		return
	}
	for i := 0; i != m.nodes.Length(); i++ {
		m.matrixCopy[i][key] = false
	}
}

func (m *graphMatrix) getFreeNodes() []Node {
	var nodes = make([]Node, 0)

	for i := 0; i != m.nodes.Length(); i++ {
		var free = true
		for j := 0; j != m.nodes.Length(); j++ {
			if m.at(i, j) {
				free = false
				break
			}
		}
		if free {
			var node, _ = m.nodes.GetByKey(i)
			nodes = append(nodes, node.(Node))
		}
	}
	return nodes
}

func (m *graphMatrix) at(i, j int) bool {
	return m.matrixCopy[i][j]
}

func getUnconnectedErrMsg(ports []Port) string {
	var nodeNames = make([]string, len(ports))
	var portTags = make([]string, len(ports))
	for i, port := range ports {
		nodeNames[i] = port.GetInnerNode().GetName()
		portTags[i] = port.GetTag()
	}
	return unconnectedPortsMsg + " nodes: " + strings.Join(nodeNames, ", ") + " tags: " + strings.Join(portTags, ", ")
}

func getUndefinedErrMsg(nodes []Node) string {
	var nodeNames = make([]string, len(nodes))
	for i, node := range nodes {
		nodeNames[i] = node.GetName()
	}
	return contextUndefinedNodes + ": " + strings.Join(nodeNames, ", ")
}

func getInaccessibleErrMsg(nodes []Node) string {
	var nodeNames = make([]string, len(nodes))
	for i, node := range nodes {
		nodeNames[i] = node.GetName()
	}
	return inaccessibleNodesMsg + ": " + strings.Join(nodeNames, ", ")
}

func sprintMatrix(m [][]bool) string {
	var result = ""
	for _, row := range m {
		var rowStr = ""
		for _, col := range row {
			if !col {
				rowStr += "_ "
			} else {
				rowStr += "+ "
			}
		}
		result += fmt.Sprintf("%s\n", rowStr)
	}
	return result
}
