package core

import "fmt"

func NewConnectivityMatrix(nodeArr []interface{}) *connectivityMatrixType {
	var nodeMap = newBiMap()
	for i, node := range nodeArr {
		nodeMap.Add(i, node)
	}

	var matrix = make([][]bool, len(nodeArr))
	for i := 0; i != len(nodeArr); i++ {
		matrix[i] = make([]bool, len(nodeArr))
	}

	return &connectivityMatrixType{
		nodes:  nodeMap,
		matrix: matrix,
	}
}

type connectivityMatrixType struct {
	nodes  *biMap
	matrix [][]bool
}

func (m *connectivityMatrixType) AddEdge(to, from interface{}) error {
	var toKey, ok1 = m.nodes.backward[to]
	if !ok1 {
		return fmt.Errorf("%v not found", to)
	}
	var fromKey, ok2 = m.nodes.backward[from]
	if !ok2 {
		return fmt.Errorf("%v not found", from)
	}
	m.matrix[toKey][fromKey] = true
	return nil
}

func (m *connectivityMatrixType) At(i, j int) bool {
	return m.at(i, j)
}

func (m *connectivityMatrixType) GetCallOrder() ([]interface{}, error) {
	var result = make([]interface{}, 0)
	var front = m.getFreeNodes()

	for len(front) != 0 {
		result = append(result, front...)
		m.removeAllDependencies(front)
		front = m.getFreeNodes()
	}

	var dependentNodes = m.getDependentNodes()
	if len(dependentNodes) > 0 {
		return nil, fmt.Errorf("nodes %v are not accessible", dependentNodes)
	}
	return result, nil
}

func (m *connectivityMatrixType) getDependentNodes() []interface{} {
	var nodes = make(map[interface{}]bool)
	for pair := range m.nodes.Iterate() {
		nodes[pair.val] = true
	}

	var freeNodes = m.getFreeNodes()
	for _, node := range freeNodes {
		delete(nodes, node)
	}

	var result = make([]interface{}, 0)
	for _, node := range nodes {
		result = append(result, node)
	}
	return result
}

func (m *connectivityMatrixType) removeAllDependencies(nodes []interface{}) {
	for _, node := range nodes {
		m.removeDependency(node)
	}
}

func (m *connectivityMatrixType) removeDependency(node interface{}) {
	var key, ok = m.nodes.GetByVal(node)
	if !ok {
		return
	}
	for i := 0; i != m.nodes.Length(); i++ {
		m.matrix[i][key] = false
	}
}

func (m *connectivityMatrixType) getFreeNodes() []interface{} {
	var nodes = make([]interface{}, 0)

	for i := 0; i != m.nodes.Length(); i++ {
		var free = true
		for j := 0; j != m.nodes.Length(); j++ {
			if m.at(i, j) {
				free = false
				break
			}
		}
		if free {
			nodes = append(nodes, m.nodes.GetByKey(i))
		}
	}
	return nodes
}

func (m *connectivityMatrixType) defaultCols() int {
	if m.defaultRows() == 0 {
		return 0
	}
	return len(m.matrix[0])
}

func (m *connectivityMatrixType) defaultRows() int {
	return len(m.matrix)
}

func (m *connectivityMatrixType) at(i, j int) bool {
	return m.matrix[j][i]
}
