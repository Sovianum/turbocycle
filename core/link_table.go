package core

import (
	"errors"
	"fmt"
)

// TODO use matrix instead of maps
type rowType map[int]bool
type linkTableType map[int]rowType

// if function fails requireTable may be in invalid networkState
func getCallOrder(requireTable, updateTable linkTableType) ([]int, error) {
	var front = getEmptyRootIds(requireTable)
	if len(front) == 0 {
		return nil, errors.New("Network has no start nodes")
	}
	deleteEmptyRoots(requireTable)

	var result = make([]int, 0)
	result = append(result, front...)

	var err error
	for len(front) != 0 {
		front, err = getNewFront(requireTable, updateTable, front)
		if err != nil {
			return nil, err
		}
		deleteEmptyRoots(requireTable)
		result = append(result, front...)

		if len(requireTable) == 0 {
			break
		}
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
	for rootId := range deleteSubTable {
		var _, ok = linkTable[rootId]
		if ok {
			for leafId := range deleteSubTable[rootId] {
				delete(linkTable[rootId], leafId)
			}
		}
	}

	return nil
}

func getSubTable(linkTable linkTableType, rowId int) linkTableType {
	var updateTable = make(linkTableType)
	updateTable[rowId] = make(rowType)

	for linkNodeId := range linkTable[rowId] {
		updateTable[rowId][linkNodeId] = true
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
