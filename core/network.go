package core

import (
	"errors"
	"fmt"
)

// TODO use matrix instead of maps
type rowType map[int]bool
type linkTableType map[int]rowType

type nodeStateType map[string]IPortState
type networkStateType map[int]nodeStateType

type Network struct {
	nodes []Node
}

func (network *Network) Solve(relaxCoef float64, maxIterNum int, precision float64) (bool, error) {
	var freePortErr = network.checkFreePorts()
	if freePortErr != nil {
		return false, freePortErr
	}

	var callOrder, callErr = getCallOrder(network.getRequireLinkTable(), network.getUpdateLinkTable())
	if callErr != nil {
		return false, callErr
	}

	var converged bool
	var err error

	for i := 0; i != maxIterNum; i++ {
		converged, err = network.makeIteration(callOrder, precision)

		if err != nil {
			err = errors.New(fmt.Sprintf(
				"Failed on iteration %d: %s", i, err.Error(),
			))
			break
		}

		if converged {
			break
		}
	}

	return converged, err
}

func (network *Network) makeIteration(callOrder []int, precision float64) (bool, error) {
	var currState, newState networkStateType
	var err error

	currState, err = network.getState()
	if err != nil {
		return false, err
	}

	newState, err = network.getNewState(callOrder)
	if err != nil {
		return false, err
	}

	var residual, residualErr = getResidual(currState, newState)
	if residualErr != nil {
		return false, residualErr
	}

	if residual <= precision {
		return true, nil
	}
	return false, nil
}

func (network *Network) getNewState(callOrder []int) (networkStateType, error) {
	for _, nodeId := range callOrder {
		var err = network.nodes[nodeId].Process()
		if err != nil {
			return nil, errors.New(fmt.Sprintf(
				"Failed on node %d: %s", nodeId, err.Error(),
			))
		}
	}
	return network.getState()
}

func (network *Network) updateNetworkState(newState networkStateType, relaxCoef float64) error {
	for nodeId, nodeState := range newState {
		if nodeId > len(network.nodes) {
			return errors.New(fmt.Sprintf("NodeId == %d > len(network.nodes) == %d", nodeId, len(network.nodes)))
		}

		for tag, portState := range nodeState {
			var port, tagErr = network.nodes[nodeId].GetPortByTag(tag)
			if tagErr != nil {
				return errors.New(fmt.Sprintf(
					"Failed to get port by tag \"%s\" from node %d: %s", tag, nodeId, tagErr.Error(),
				))
			}

			var newPortState, stateErr = port.GetState().Mix(portState, relaxCoef)
			if stateErr != nil {
				return errors.New(fmt.Sprintf(
					"Failed to mix state of port \"%s\" from node %d: %s", tag, nodeId, stateErr.Error(),
				))
			}

			port.SetState(newPortState)
		}
	}

	return nil
}

func (network *Network) getState() (networkStateType, error) {
	var result = make(networkStateType)

	for nodeId, node := range network.nodes {
		var nodeState = make(nodeStateType)

		var tags = node.GetPortTags()
		for _, tag := range tags {
			var port, err = node.GetPortByTag(tag)
			if err != nil {
				return nil, err
			}
			nodeState[tag] = port.GetState()
		}
		result[nodeId] = nodeState
	}

	return result, nil
}

func (network *Network) checkFreePorts() error {
	for nodeId, node := range network.nodes {
		for portTag, port := range node.GetPorts() {
			if port.linkPort == nil {
				return errors.New(fmt.Sprintf("Found free port \"%s\" of node %d", portTag, nodeId))
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

func (network *Network) getLinkTable(leafExtractor func(Node) []Node) linkTableType {
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

func getResidual(state1, state2 networkStateType) (float64, error) {
	var result float64 = 0

	for nodeKey, nodeState1 := range state1 {
		var nodeState2, ok = state2[nodeKey]
		if !ok {
			return 0, errors.New(fmt.Sprintf("Node %d not found in state2", nodeKey))
		}

		if len(nodeState1) != len(nodeState2) {
			return 0, errors.New(fmt.Sprintf(
				"States of node %d has different lengths (%d, %d)", nodeKey, len(nodeState1), len(nodeState2)),
			)
		}

		for portKey := range nodeState1 {
			var residual, err = nodeState1[portKey].MaxResidual(nodeState2[portKey])
			if err != nil {
				return 0, errors.New(fmt.Sprintf(
					"Failed to get residual of node %d at port %d: %s", nodeKey, portKey, err.Error(),
				))
			}
			if residual > result {
				result = residual
			}
		}
	}

	return result, nil
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
