package parametric

import (
	"github.com/Sovianum/turbocycle/core/graph"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/material/gases"
)

type Efficient interface {
	Efficiency() float64
}

func NewGasPart(gas gases.Gas, tStagIn, pStagIn, pStagOut float64) *GasPart {
	return &GasPart{
		GasSource:            source.NewGasSourceNode(gas),
		TemperatureSource:    source.NewTemperatureSourceNode(tStagIn),
		InputPressureSource:  source.NewPressureSourceNode(pStagIn),
		OutputPressureSource: source.NewPressureSourceNode(pStagOut),
	}
}

type GasPart struct {
	GasSource            source.GasSourceNode
	TemperatureSource    source.TemperatureSourceNode
	InputPressureSource  source.PressureSourceNode
	OutputPressureSource source.PressureSourceNode
}

func (part *GasPart) Nodes() []graph.Node {
	return []graph.Node{
		part.GasSource, part.TemperatureSource,
		part.InputPressureSource, part.OutputPressureSource,
	}
}

func NewGasGeneratorPart(
	compressor c.ParametricCompressorNode,
	burner c.ParametricBurnerNode,
	turbine c.ParametricTurbineNode,
	shaft c.TransmissionNode,
	compressorPipe c.PressureLossNode,
) *GasGeneratorPart {
	var result = &GasGeneratorPart{
		TurboShaftPart: NewTurboShaftPart(compressor, turbine, shaft),
		Burner:         burner,
		CompressorPipe: compressorPipe,
	}
	graph.LinkAll(
		[]graph.Port{
			compressor.GasOutput(), compressor.TemperatureOutput(),
			compressor.PressureOutput(), compressor.MassRateOutput(),
		},
		[]graph.Port{
			compressorPipe.GasInput(), compressorPipe.TemperatureInput(),
			compressorPipe.PressureInput(), compressorPipe.MassRateInput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			compressorPipe.GasOutput(), compressorPipe.TemperatureOutput(),
			compressorPipe.PressureOutput(), compressorPipe.MassRateOutput(),
		},
		[]graph.Port{
			burner.GasInput(), burner.TemperatureInput(),
			burner.PressureInput(), burner.MassRateInput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			burner.GasOutput(), burner.TemperatureOutput(),
			burner.PressureOutput(),
		},
		[]graph.Port{
			turbine.GasInput(), turbine.TemperatureInput(),
			turbine.PressureInput(),
		},
	)
	sink.SinkPort(burner.MassRateOutput())
	sink.SinkAll(burner.MassRateOutput(), turbine.MassRateInput())
	return result
}

type GasGeneratorPart struct {
	*TurboShaftPart
	Burner         c.ParametricBurnerNode
	CompressorPipe c.PressureLossNode
}

func (part *GasGeneratorPart) Nodes() []graph.Node {
	return append(
		part.TurboShaftPart.Nodes(),
		part.Burner, part.CompressorPipe,
	)
}

func NewTurboShaftPart(
	compressor c.ParametricCompressorNode,
	turbine c.ParametricTurbineNode, shaft c.TransmissionNode,
) *TurboShaftPart {
	var result = &TurboShaftPart{
		Compressor: compressor,
		Turbine:    turbine,
		Shaft:      shaft,
	}
	graph.Link(result.Compressor.PowerOutput(), result.Shaft.PowerInput())
	graph.Link(result.Compressor.RPMOutput(), result.Turbine.RPMInput())
	sink.SinkAll(shaft.PowerOutput(), turbine.PowerOutput())
	return result
}

type TurboShaftPart struct {
	Compressor c.ParametricCompressorNode
	Turbine    c.ParametricTurbineNode
	Shaft      c.TransmissionNode
}

func (part *TurboShaftPart) Nodes() []graph.Node {
	return []graph.Node{
		part.Compressor, part.Turbine, part.Shaft,
	}
}
