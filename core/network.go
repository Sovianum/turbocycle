package core

import (
	"encoding/json"
	"errors"
	"fmt"
)

type nodeStateType map[string]PortState
type networkStateType map[string]nodeStateType

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

type Network interface {
	json.Marshaler
	Solve(relaxCoef float64, maxIterNum int, precision float64) (bool, error)
	GetCallOrder() ([]string, error)
	GetState() (networkStateType, error)
}

type network struct {
	nodes map[string]Node
}

func NewNetwork(nodes map[string]Node) Network {
	return &network{nodes: nodes}
}

func (network *network) MarshalJSON() ([]byte, error) {
	return json.Marshal(network.nodes)
}

func (network *network) Solve(relaxCoef float64, maxIterNum int, precision float64) (bool, error) {
	var callOrder, callErr = network.getCallOrder()
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

func (network *network) GetCallOrder() ([]string, error) {
	return network.getCallOrder()
}

func (network *network) GetState() (networkStateType, error) {
	return network.getState()
}

func (network *network) getCallOrder() ([]string, error) {
	var freePortErr = network.checkFreePorts()
	if freePortErr != nil {
		return nil, freePortErr
	}

	var contextDefinitionErr = network.checkContextDefinition()
	if contextDefinitionErr != nil {
		return nil, contextDefinitionErr
	}

	var requireLinkTable, requireTableErr = network.getRequireLinkTable()
	if requireTableErr != nil {
		return nil, requireTableErr
	}

	var updateLinkTable, updateTableErr = network.getUpdateLinkTable()
	if updateTableErr != nil {
		return nil, updateTableErr
	}

	var callOrder, callErr = getCallOrder(requireLinkTable, updateLinkTable)
	return callOrder, callErr
}

func (network *network) makeIteration(callOrder []string, precision float64) (bool, error) {
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

func (network *network) checkContextDefinition() error {
	var nodeKeys = make([]string, 0)
	for key, node := range network.nodes {
		if !node.ContextDefined() {
			nodeKeys = append(nodeKeys, key)
		}
	}

	if len(nodeKeys) > 0 {
		return errors.New(fmt.Sprintf("Nodes %v are not context defined", nodeKeys))
	}
	return nil
}

func (network *network) getNewState(callOrder []string) (networkStateType, error) {
	for _, nodeKey := range callOrder {
		var err = network.nodes[nodeKey].Process()
		if err != nil {
			return nil, errors.New(fmt.Sprintf(
				"Failed on node %d: %s", nodeKey, err.Error(),
			))
		}
	}
	return network.getState()
}

func (network *network) updateNetworkState(newState networkStateType, relaxCoef float64) error {
	for nodeKey, nodeState := range newState {
		for tag, portState := range nodeState {
			var port, tagErr = network.nodes[nodeKey].GetPortByTag(tag)
			if tagErr != nil {
				return errors.New(fmt.Sprintf(
					"Failed to get portType by tag \"%s\" from node %d: %s", tag, nodeKey, tagErr.Error(),
				))
			}

			var newPortState, stateErr = port.GetState().Mix(portState, relaxCoef)
			if stateErr != nil {
				return errors.New(fmt.Sprintf(
					"Failed to mix networkState of portType \"%s\" from node %d: %s", tag, nodeKey, stateErr.Error(),
				))
			}

			port.SetState(newPortState)
		}
	}

	return nil
}

func (network *network) getState() (networkStateType, error) {
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

func (network *network) checkFreePorts() error {
	for nodeId, node := range network.nodes {
		for portTag, port := range node.GetPorts() {
			if port.GetLinkPort() == nil {
				return errors.New(fmt.Sprintf("Found free portType \"%s\" of node %d", portTag, nodeId))
			}
		}
	}
	return nil
}

func (network *network) getUpdateLinkTable() (linkTableType, error) {
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

func (network *network) getRequireLinkTable() (linkTableType, error) {
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

func (network *network) getLinkTable(leafExtractor func(Node) ([]Node, error)) (linkTableType, error) {
	var idMap = make(map[Node]string)
	for key, node := range network.nodes {
		idMap[node] = key
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
