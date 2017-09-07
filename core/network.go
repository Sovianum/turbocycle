package core

import (
	"errors"
	"fmt"
)

// TODO use matrix instead of maps
type rowType map[int]bool
type linkTableType map[int]rowType

type Network struct {
	nodes []Node
}

func (network *Network) checkFreePorts() error {
	for _, node := range network.nodes {
		for _, port := range node.GetPorts() {
			if port.linkPort == nil {
				return errors.New("Found free port")
			}
		}
	}
	return nil
}

func (network *Network) getUpdateLinkTable() linkTableType {
	return network.getLinkTable(func(node Node) []Node {
		var result = make([]Node, 0)
		for _, portTag := range node.GetUpdatePortTags() {
			result = append(result, node.GetPorts()[portTag].GetOuterNode())
		}
		return result
	})
}

func (network *Network) getRequireLinkTable() linkTableType {
	return network.getLinkTable(func(node Node) []Node {
		var result = make([]Node, 0)
		for _, portTag := range node.GetRequirePortTags() {
			result = append(result, node.GetPorts()[portTag].GetOuterNode())
		}
		return result
	})
}

func (network *Network) getLinkTable(leafExtractor func(Node)[]Node) linkTableType {
	var idMap = make(map[Node]int)
	for i, node := range network.nodes {
		idMap[node] = i
	}

	var result = make(linkTableType)
	for _, root := range network.nodes {
		result[idMap[root]] = make(rowType)
		for _, leaf := range leafExtractor(root) {
			result[idMap[root]][idMap[leaf]] = true
		}
	}
	return result
}

// if function fails requireTable may be in invalid state
func getCallOrder(requireTable, updateTable linkTableType) ([]int, error) {
	var front = getEmptyRootIds(requireTable)
	if len(front) == 0 {
		return nil, errors.New("Network has not start nodes")
	}
	deleteEmptyRoots(requireTable)

	var result = make([]int, 0)
	result = append(result, front...)

	var err error
	for len(front) != 0 {
		if len(requireTable) == 0 {
			break
		}

		front, err = getNewFront(requireTable, updateTable, front)
		if err != nil {
			return nil, err
		}
		deleteEmptyRoots(requireTable)
		result = append(result, front...)
	}

	if len(requireTable) != 0 {
		var inaccessibleRoots = make([]int, 0)
		for rootId := range requireTable {
			inaccessibleRoots = append(inaccessibleRoots, rootId)
		}

		return nil, errors.New(fmt.Sprintf("Roots %v can not be called", inaccessibleRoots))
	}

	return result, nil
}

func getNewFront(requireTable, updateTable linkTableType, frontIds []int) ([]int, error) {
	for _, id := range frontIds {
		var subTable = getSubTable(updateTable, id)
		var invertedSubTable = invertTable(subTable)

		var err = subtractSubTable(requireTable, invertedSubTable)
		if err != nil {
			return nil, err
		}
	}
	return getEmptyRootIds(requireTable), nil
}

func subtractSubTable(linkTable, deleteSubTable linkTableType) error {
	var err = isSubTable(linkTable, deleteSubTable)
	if err != nil {
		return err
	}

	for rootId := range deleteSubTable {
		for leafId := range deleteSubTable[rootId] {
			delete(linkTable[rootId], leafId)
		}
	}

	return nil
}

func isSubTable(table, subTable linkTableType) error {
	for rootId := range subTable {
		var _, ok = table[rootId]
		if !ok {
			return errors.New(fmt.Sprintf("Root with id = %d not found", rootId))
		}

		for leafId := range subTable[rootId] {
			var _, ok = table[rootId][leafId]
			if !ok {
				return errors.New(fmt.Sprintf("Leaf with id = %d of root with id = %d not found", leafId, rootId))
			}
		}
	}
	return nil
}

func getSubTable(linkTable linkTableType, rowId int) linkTableType {
	var updateTable = make(linkTableType)
	updateTable[rowId] = make(rowType)

	for leafNodeId := range linkTable[rowId] {
		updateTable[rowId][leafNodeId] = true
	}

	return updateTable
}

func deleteEmptyRoots(table linkTableType) {
	for rootId := range table {
		if len(table[rootId]) == 0 {
			delete(table, rootId)
		}
	}
}

func getEmptyRootIds(table linkTableType) []int {
	var result = make([]int, 0)
	for rootId, linkMap := range table {
		if len(linkMap) == 0 {
			result = append(result, rootId)
		}
	}
	return result
}

func invertTable(table linkTableType) linkTableType {
	var result = make(linkTableType)
	for rootId := range table {
		for leafId := range table[rootId] {
			var _, ok = result[leafId]
			if !ok {
				result[leafId] = make(rowType)
			}

			result[leafId][rootId] = true
		}
	}

	return result
}
