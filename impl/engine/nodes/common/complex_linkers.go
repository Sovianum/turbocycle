package common

import "github.com/Sovianum/turbocycle/core/graph"

func LinkComplexInToOut(node1 ComplexGasSink, node2 ComplexGasSource) {
	LinkComplexOutToIn(node2, node1)
}

func LinkComplexOutToIn(node1 ComplexGasSource, node2 ComplexGasSink) {
	graph.LinkAll(
		[]graph.Port{node1.GasOutput(), node1.TemperatureOutput(), node1.PressureOutput(), node1.MassRateOutput()},
		[]graph.Port{node2.GasInput(), node2.TemperatureInput(), node2.PressureInput(), node2.MassRateInput()},
	)
}

func LinkComplexOutToOut(node1 ComplexGasSource, node2 ComplexGasSource) {
	graph.LinkAll(
		[]graph.Port{node1.GasOutput(), node1.TemperatureOutput(), node1.PressureOutput(), node1.MassRateOutput()},
		[]graph.Port{node2.GasOutput(), node2.TemperatureOutput(), node2.PressureOutput(), node2.MassRateOutput()},
	)
}

func LinkComplexInToIn(node1 ComplexGasSink, node2 ComplexGasSink) {
	graph.LinkAll(
		[]graph.Port{node1.GasInput(), node1.TemperatureInput(), node1.PressureInput(), node1.MassRateInput()},
		[]graph.Port{node2.GasInput(), node2.TemperatureInput(), node2.PressureInput(), node2.MassRateInput()},
	)
}
