package helper

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
)

func NewWeakPseudoComplexGasSource(source nodes.ComplexGasSource) nodes.ComplexGasSource {
	return NewPseudoComplexGasSource(
		graph.NewWeakPort(source.GasOutput()),
		graph.NewWeakPort(source.TemperatureOutput()),
		graph.NewWeakPort(source.PressureOutput()),
		graph.NewWeakPort(source.MassRateOutput()),
	)
}

func NewPseudoComplexGasSource(
	gasOutput, temperatureOutput,
	pressureOutput, massRateOutput graph.Port,
) nodes.ComplexGasSource {
	return &pseudoComplexGasSource{
		gasOutput:         gasOutput,
		temperatureOutput: temperatureOutput,
		pressureOutput:    pressureOutput,
		massRateOutput:    massRateOutput,
	}
}

type pseudoComplexGasSource struct {
	temperatureOutput graph.Port
	pressureOutput    graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port
}

func (s *pseudoComplexGasSource) GasOutput() graph.Port {
	return s.gasOutput
}

func (s *pseudoComplexGasSource) TemperatureOutput() graph.Port {
	return s.temperatureOutput
}

func (s *pseudoComplexGasSource) PressureOutput() graph.Port {
	return s.pressureOutput
}

func (s *pseudoComplexGasSource) MassRateOutput() graph.Port {
	return s.massRateOutput
}

func NewWeakPseudoComplexGasSink(sink nodes.ComplexGasSink) nodes.ComplexGasSink {
	return NewPseudoComplexGasSink(
		graph.NewWeakPort(sink.GasInput()),
		graph.NewWeakPort(sink.TemperatureInput()),
		graph.NewWeakPort(sink.PressureInput()),
		graph.NewWeakPort(sink.MassRateInput()),
	)
}

func NewPseudoComplexGasSink(
	gasInput, temperatureInput,
	pressureInput, massRateInput graph.Port,
) nodes.ComplexGasSink {
	return &pseudoComplexGasSink{
		gasInput:         gasInput,
		temperatureInput: temperatureInput,
		pressureInput:    pressureInput,
		massRateInput:    massRateInput,
	}
}

type pseudoComplexGasSink struct {
	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port
	massRateInput    graph.Port
}

func (s *pseudoComplexGasSink) GasInput() graph.Port {
	return s.gasInput
}

func (s *pseudoComplexGasSink) TemperatureInput() graph.Port {
	return s.temperatureInput
}

func (s *pseudoComplexGasSink) PressureInput() graph.Port {
	return s.pressureInput
}

func (s *pseudoComplexGasSink) MassRateInput() graph.Port {
	return s.massRateInput
}
