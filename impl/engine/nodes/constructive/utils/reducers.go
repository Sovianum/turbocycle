package utils

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

func NewEquality(port1, port2 graph.Port) graph.ReduceNode {
	var groupReduceFunc = func(curr float64, port graph.Port) (float64, error) {
		var val1 = port1.GetState().Value()
		var val2 = port2.GetState().Value()

		var floatVal1, ok1 = val1.(float64)
		var floatVal2, ok2 = val2.(float64)

		if !ok1 {
			return 0, common.GetTypeError("float64", val1)
		}

		if !ok2 {
			return 0, common.GetTypeError("float64", val2)
		}

		return floatVal1 - floatVal2, nil
	}

	var totalReduceFunc = func(curr, new float64) (float64, error) {
		return new, nil
	}

	var result = graph.NewReduceNode(
		groupReduceFunc, totalReduceFunc, 1, 0,
	)
	result.AddPortGroup(port1, port2)
	return result
}

func NewSum(ports ...graph.Port) graph.ReduceNode {
	var result = NewMultiAdder()
	for _, port := range ports {
		result.AddPortGroup(port)
	}
	return result
}

func NewMultiAdderFromPorts(ports ...[]graph.Port) graph.ReduceNode {
	var result = NewMultiAdder()
	for _, group := range ports {
		result.AddPortGroup(group...)
	}
	return result
}

func NewMultiAdder() graph.ReduceNode {
	groupReduceFunc := func(curr float64, port graph.Port) (float64, error) {
		val := port.GetState().Value()
		floatVal, ok := val.(float64)

		if !ok {
			return 0, common.GetTypeError("float64", val)
		}

		return curr * floatVal, nil
	}
	totalReduceFunc := func(curr, new float64) (float64, error) {
		return curr + new, nil
	}
	return graph.NewReduceNode(
		groupReduceFunc, totalReduceFunc, 1, 0,
	)
}
