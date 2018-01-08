package constructive

import (
	"math"

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

func NewMultiAdderFromPorts(ports [][]graph.Port) graph.ReduceNode {
	var result = NewMultiAdder()
	for _, group := range ports {
		result.AddPortGroup(group...)
	}
	return result
}

func NewMultiAdder() graph.ReduceNode {
	var groupReduceFunc = func(curr float64, port graph.Port) (float64, error) {
		var val = port.GetState().Value()
		var floatVal, ok = val.(float64)

		if !ok {
			return 0, common.GetTypeError("float64", val)
		}

		return curr * floatVal, nil
	}

	var totalReduceFunc = func(curr, new float64) (float64, error) {
		return curr + new, nil
	}

	return graph.NewReduceNode(
		groupReduceFunc, totalReduceFunc, 1, 0,
	)
}

func NormMassRate(massRate, massRate0, t, t0, p, p0 float64) float64 {
	return massRate / massRate0 * math.Sqrt(t0/t) * p / p0
}

func MassRate(normMassRate, massRate0, t, t0, p, p0 float64) float64 {
	return normMassRate / (math.Sqrt(t0/t) * p / p0) * massRate0
}

func NormRpm(rpm, rpm0, t, t0 float64) float64 {
	return rpm / rpm0 * math.Sqrt(t0/t)
}

func Rpm(normRpm, rpm0, t, t0 float64) float64 {
	return normRpm / math.Sqrt(t0/t) * rpm0
}
