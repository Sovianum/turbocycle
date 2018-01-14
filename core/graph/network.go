package graph

import (
	"fmt"
)

type nodeStateType map[Port]PortState
type networkStateType map[Node]nodeStateType

type Network interface {
	// Solve solves the graph and returns pair (isConverged bool and solutionErr error)
	// before skipIterations is reached residuals are not checked. This is done to eliminate
	// necessity to initialize all the ports before Solve call
	Solve(relaxCoef float64, skipIterations int, maxIterNum int, precision float64) (bool, error)
}

func NewNetwork(nodes []Node) (Network, GraphError) {
	var matrix, err = newGraphMatrix(nodes)
	if err != nil {
		return nil, err
	}

	return &network{nodes: nodes, graphMatrix: matrix}, nil
}

type network struct {
	nodes       []Node
	graphMatrix *graphMatrix
}

func (network *network) Solve(relaxCoef float64, skipIterations int, maxIterNum int, precision float64) (bool, error) {
	var callOrder, callErr = network.graphMatrix.GetCallOrder()
	if callErr != nil {
		return false, callErr
	}

	var converged bool
	var err error

	for i := 0; i != maxIterNum; i++ {
		converged, err = network.getStates(callOrder, i >= skipIterations, precision)

		if err != nil {
			err = fmt.Errorf(
				"failed on iteration %d: %s", i, err.Error(),
			)
			break
		}

		if converged {
			break
		}
	}

	return converged, err
}

func (network *network) getStates(callOrder []Node, needCheck bool, precision float64) (bool, error) {
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

	if !needCheck {
		return false, nil
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

func (network *network) getNewState(callOrder []Node) (networkStateType, error) {
	for _, node := range callOrder {
		var err = node.Process()
		if err != nil {
			return nil, fmt.Errorf(
				"failed on node with name %s: %s", node.GetName(), err.Error(),
			)
		}
	}
	return network.getState()
}

func (network *network) getState() (networkStateType, error) {
	var result = make(networkStateType)

	for _, node := range network.nodes {
		var nodeState = make(nodeStateType)

		var ports = node.GetPorts()
		for _, port := range ports {
			nodeState[port] = port.GetState()
		}
		result[node] = nodeState
	}

	return result, nil
}

func getResidual(state1, state2 networkStateType) (float64, error) {
	var result float64 = 0

	for node, nodeState1 := range state1 {
		var nodeState2, ok = state2[node]
		if !ok {
			return 0, fmt.Errorf("node %v with name %s not found in state2", node, node.GetName())
		}

		if len(nodeState1) != len(nodeState2) {
			return 0, fmt.Errorf(
				"states of node %v with name %s has different lengths (%d, %d)",
				node, node.GetName(), len(nodeState1), len(nodeState2),
			)
		}

		for port := range nodeState1 {
			var portState1 = nodeState1[port]
			var portState2 = nodeState2[port]

			if portState1 == nil && portState2 == nil {
				continue
			}

			if portState1 == nil {
				return 0, fmt.Errorf(
					"port with tag %s of node %s has nil state on curr step",
					port.GetTag(), port.GetInnerNode().GetName(),
				)
			}
			if portState2 == nil {
				return 0, fmt.Errorf(
					"port with tag %s of node %s has nil state on new step",
					port.GetTag(), port.GetInnerNode().GetName(),
				)
			}

			var residual, err = portState1.MaxResidual(portState2)
			if err != nil {
				return 0, fmt.Errorf(
					"failed to get residual of node %v with name %s at portType %s: %s",
					node, node.GetName(), port, err.Error(),
				)
			}
			if residual > result {
				result = residual
			}
		}
	}

	return result, nil
}
