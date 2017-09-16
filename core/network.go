package core

import (
	"errors"
	"fmt"
)

type nodeStateType map[string]PortState
type networkStateType map[int]nodeStateType

type Network struct {
	nodes []Node
}

func (networkState networkStateType) String() string {
	var result = ""
	for key, nodeState := range networkState {
		result += fmt.Sprintf("%d\n", key)
		for tag, val := range nodeState {
			result += fmt.Sprintf("\t%s\t%v\n", tag, val)
		}
	}

	return result
}

func NewNetwork(nodes []Node) *Network {
	return &Network{nodes: nodes}
}

func (network *Network) Solve(relaxCoef float64, maxIterNum int, precision float64) (bool, error) {
	var freePortErr = network.checkFreePorts()
	if freePortErr != nil {
		return false, freePortErr
	}

	var contextDefinitionErr = network.checkContextDefinition()
	if contextDefinitionErr != nil {
		return false, contextDefinitionErr
	}

	var requireLinkTable, requireTableErr = network.getRequireLinkTable()
	if requireTableErr != nil {
		return false, requireTableErr
	}

	var updateLinkTable, updateTableErr = network.getUpdateLinkTable()
	if updateTableErr != nil {
		return false, updateTableErr
	}

	var callOrder, callErr = getCallOrder(requireLinkTable, updateLinkTable)
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

func (network *Network) GetState() (networkStateType, error) {
	return network.getState()
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

func (network *Network) checkContextDefinition() error {
	var nodeIds = make([]int, 0)
	for id, node := range network.nodes {
		if !node.ContextDefined() {
			nodeIds = append(nodeIds, id)
		}
	}

	if len(nodeIds) > 0 {
		return errors.New(fmt.Sprintf("Nodes %v are not context defined", nodeIds))
	}
	return nil
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
					"Failed to get portType by tag \"%s\" from node %d: %s", tag, nodeId, tagErr.Error(),
				))
			}

			var newPortState, stateErr = port.GetState().Mix(portState, relaxCoef)
			if stateErr != nil {
				return errors.New(fmt.Sprintf(
					"Failed to mix networkState of portType \"%s\" from node %d: %s", tag, nodeId, stateErr.Error(),
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
			if port.GetLinkPort() == nil {
				return errors.New(fmt.Sprintf("Found free portType \"%s\" of node %d", portTag, nodeId))
			}
		}
	}
	return nil
}

func (network *Network) getUpdateLinkTable() (linkTableType, error) {
	return network.getLinkTable(func(node Node) ([]Node, error) {
		var result = make([]Node, 0)

		var updatePorts, err = node.GetUpdatePortTags()
		if err != nil {
			return nil, err
		}

		for _, portTag := range updatePorts {
			result = append(result, node.GetPorts()[portTag].GetOuterNode())
		}
		return result, nil
	})
}

func (network *Network) getRequireLinkTable() (linkTableType, error) {
	return network.getLinkTable(func(node Node) ([]Node, error) {
		var result = make([]Node, 0)

		var requirePorts, err = node.GetRequirePortTags()
		if err != nil {
			return nil, err
		}

		for _, portTag := range requirePorts {
			result = append(result, node.GetPorts()[portTag].GetOuterNode())
		}
		return result, nil
	})
}

func (network *Network) getLinkTable(leafExtractor func(Node) ([]Node, error)) (linkTableType, error) {
	var idMap = make(map[Node]int)
	for i, node := range network.nodes {
		idMap[node] = i
	}

	var result = make(linkTableType)
	var nodes []Node
	var err error

	for _, root := range network.nodes {
		result[idMap[root]] = make(rowType)

		nodes, err = leafExtractor(root)
		if err != nil {
			break
		}
		for _, leaf := range nodes {
			result[idMap[root]][idMap[leaf]] = true
		}
	}

	if err != nil {
		return nil, err
	}
	return result, nil
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
					"Failed to get residual of node %d at portType %d: %s", nodeKey, portKey, err.Error(),
				))
			}
			if residual > result {
				result = residual
			}
		}
	}

	return result, nil
}
