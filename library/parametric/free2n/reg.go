package free2n

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/variator"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive/utils"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/library/parametric"
	"github.com/Sovianum/turbocycle/material/gases"
)

type DoubleShaftRegFreeScheme interface {
	DoubleShaftFreeScheme
	Regenerator() c.RegeneratorNode
}

func NewDoubleShaftRegFreeScheme(
	gas gases.Gas, tAtm, pAtmIn, pAtmOut,
	tGas, etaM float64,
	compressor c.ParametricCompressorNode,
	compressorPipe c.PressureLossNode,
	regenerator c.RegeneratorNode,
	breaker helper.ComplexCycleBreakNode,
	burner c.ParametricBurnerNode,
	compressorTurbine c.ParametricTurbineNode,
	compressorTurbinePipe c.PressureLossNode,
	freeTurbine c.ParametricTurbineNode,
	freeTurbinePipe c.PressureLossNode,
	payload c.Payload,
) DoubleShaftRegFreeScheme {
	result := &doubleShaftRegFreeScheme{
		gasPart: parametric.NewGasPart(gas, tAtm, pAtmIn, pAtmOut),
		gasGeneratorPart: parametric.NewGasGeneratorPart(
			compressor, burner, compressorTurbine, c.NewTransmissionNode(etaM), compressorPipe,
		),
		regenerator: regenerator,
		breaker:     breaker,

		ctPipe:   compressorTurbinePipe,
		fTurbine: freeTurbine,
		ftPipe:   freeTurbinePipe,
		payload:  payload,

		burnerTemperatureSource: source.NewTemperatureSourceNode(tGas),

		assembler: graph.NewVectorAssemblerNode(),

		variators: []variator.Variator{
			variator.FromCallables(compressor.NormMassRate, compressor.SetNormMassRate),
			variator.FromCallables(compressor.NormPiStag, compressor.SetNormPiStag),
			variator.FromCallables(burner.FuelRateRel, burner.SetFuelRateRel),
			variator.FromCallables(payload.NormRPM, payload.SetNormRPM),
			variator.FromCallables(compressorTurbine.NormPiT, compressorTurbine.SetNormPiT),
			variator.FromCallables(freeTurbine.NormPiT, freeTurbine.SetNormPiT),
		},
	}
	result.linkPorts()
	result.setEquations()
	return result
}

type doubleShaftRegFreeScheme struct {
	gasPart          *parametric.GasPart
	gasGeneratorPart *parametric.GasGeneratorPart
	regenerator      c.RegeneratorNode
	breaker          helper.ComplexCycleBreakNode

	ctPipe   c.PressureLossNode
	fTurbine c.ParametricTurbineNode
	ftPipe   c.PressureLossNode
	payload  c.Payload

	burnerTemperatureSource source.TemperatureSourceNode

	gasGenMassRateEq      graph.ReduceNode
	gasGenPowerEq         graph.ReduceNode
	freeTurbineMassRateEq graph.ReduceNode
	freeTurbinePowerEq    graph.ReduceNode
	freeTurbinePressureEq graph.ReduceNode
	gasGenBurnerEq        graph.ReduceNode

	assembler graph.VectorAssemblerNode
	variators []variator.Variator
}

func (scheme *doubleShaftRegFreeScheme) TemperatureSource() source.TemperatureSourceNode {
	return scheme.burnerTemperatureSource
}

func (scheme *doubleShaftRegFreeScheme) Regenerator() c.RegeneratorNode {
	return scheme.regenerator
}

func (scheme *doubleShaftRegFreeScheme) Compressor() c.ParametricCompressorNode {
	return scheme.gasGeneratorPart.Compressor
}

func (scheme *doubleShaftRegFreeScheme) CompressorPipe() c.PressureLossNode {
	return scheme.gasGeneratorPart.CompressorPipe
}

func (scheme *doubleShaftRegFreeScheme) Burner() c.ParametricBurnerNode {
	return scheme.gasGeneratorPart.Burner
}

func (scheme *doubleShaftRegFreeScheme) CompressorTurbine() c.ParametricTurbineNode {
	return scheme.gasGeneratorPart.Turbine
}

func (scheme *doubleShaftRegFreeScheme) CompressorTurbinePipe() c.PressureLossNode {
	return scheme.ctPipe
}

func (scheme *doubleShaftRegFreeScheme) FreeTurbine() c.ParametricTurbineNode {
	return scheme.fTurbine
}

func (scheme *doubleShaftRegFreeScheme) FreeTurbinePipe() c.PressureLossNode {
	return scheme.ftPipe
}

func (scheme *doubleShaftRegFreeScheme) Payload() c.Payload {
	return scheme.payload
}

func (scheme *doubleShaftRegFreeScheme) Assembler() graph.VectorAssemblerNode {
	return scheme.assembler
}

func (scheme *doubleShaftRegFreeScheme) Variators() []variator.Variator {
	return scheme.variators
}

func (scheme *doubleShaftRegFreeScheme) GetNetwork() (graph.Network, error) {
	var nodes = append(scheme.gasPart.Nodes(), scheme.gasGeneratorPart.Nodes()...)
	nodes = append(
		nodes, scheme.ctPipe, scheme.ftPipe, scheme.fTurbine, scheme.payload, scheme.burnerTemperatureSource,
		scheme.assembler, scheme.regenerator, scheme.breaker, scheme.gasGenMassRateEq, scheme.gasGenPowerEq, scheme.freeTurbineMassRateEq,
		scheme.freeTurbinePowerEq, scheme.freeTurbinePressureEq, scheme.gasGenBurnerEq,
	)
	return graph.NewNetwork(nodes)
}

func (scheme *doubleShaftRegFreeScheme) linkPorts() {
	graph.LinkAll(
		[]graph.Port{
			scheme.gasPart.GasSource.GasOutput(),
			scheme.gasPart.TemperatureSource.TemperatureOutput(),
			scheme.gasPart.InputPressureSource.PressureOutput(),
			scheme.gasPart.OutputPressureSource.PressureOutput(),
		},
		[]graph.Port{
			scheme.gasGeneratorPart.Compressor.GasInput(),
			scheme.gasGeneratorPart.Compressor.TemperatureInput(),
			scheme.gasGeneratorPart.Compressor.PressureInput(),
			scheme.ftPipe.PressureOutput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			scheme.gasGeneratorPart.Compressor.GasOutput(),
			scheme.gasGeneratorPart.Compressor.TemperatureOutput(),
			scheme.gasGeneratorPart.Compressor.PressureOutput(),
			scheme.gasGeneratorPart.Compressor.MassRateOutput(),
		},
		[]graph.Port{
			scheme.regenerator.ColdInput().GasInput(),
			scheme.regenerator.ColdInput().TemperatureInput(),
			scheme.regenerator.ColdInput().PressureInput(),
			scheme.regenerator.ColdInput().MassRateInput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			scheme.regenerator.ColdOutput().GasOutput(),
			scheme.regenerator.ColdOutput().TemperatureOutput(),
			scheme.regenerator.ColdOutput().PressureOutput(),
			scheme.regenerator.ColdOutput().MassRateOutput(),
		},
		[]graph.Port{
			scheme.gasGeneratorPart.CompressorPipe.GasInput(),
			scheme.gasGeneratorPart.CompressorPipe.TemperatureInput(),
			scheme.gasGeneratorPart.CompressorPipe.PressureInput(),
			scheme.gasGeneratorPart.CompressorPipe.MassRateInput(),
		},
	)
	sink.SinkAll(
		scheme.gasGeneratorPart.Compressor.MassRateInput(),
		scheme.ftPipe.PressureInput(),
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.gasGeneratorPart.Turbine.GasOutput(),
			scheme.gasGeneratorPart.Turbine.TemperatureOutput(),
			scheme.gasGeneratorPart.Turbine.PressureOutput(),
			scheme.gasGeneratorPart.Turbine.MassRateOutput(),
		},
		[]graph.Port{
			scheme.ctPipe.GasInput(),
			scheme.ctPipe.TemperatureInput(),
			scheme.ctPipe.PressureInput(),
			scheme.ctPipe.MassRateInput(),
		},
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.ctPipe.GasOutput(),
			scheme.ctPipe.TemperatureOutput(),
			scheme.ctPipe.PressureOutput(),
		},
		[]graph.Port{
			scheme.fTurbine.GasInput(),
			scheme.fTurbine.TemperatureInput(),
			scheme.fTurbine.PressureInput(),
		},
	)
	sink.SinkAll(
		scheme.ctPipe.MassRateOutput(), scheme.fTurbine.MassRateInput(),
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.fTurbine.GasOutput(), scheme.fTurbine.TemperatureOutput(),
			scheme.fTurbine.MassRateOutput(),
		},
		[]graph.Port{
			scheme.ftPipe.GasInput(), scheme.ftPipe.TemperatureInput(),
			scheme.ftPipe.MassRateInput(),
		},
	)
	sink.SinkAll(scheme.fTurbine.PressureOutput())
	graph.LinkAll(
		[]graph.Port{
			scheme.ftPipe.GasOutput(), scheme.ftPipe.TemperatureOutput(),
			graph.NewWeakPort(scheme.ftPipe.PressureOutput()), scheme.ftPipe.MassRateOutput(),
		},
		[]graph.Port{
			scheme.breaker.GasInput(), scheme.breaker.TemperatureInput(),
			scheme.breaker.PressureInput(), scheme.breaker.MassRateInput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			scheme.breaker.GasOutput(), scheme.breaker.TemperatureOutput(),
			scheme.breaker.PressureOutput(), scheme.breaker.MassRateOutput(),
		},
		[]graph.Port{
			scheme.regenerator.HotInput().GasInput(), scheme.regenerator.HotInput().TemperatureInput(),
			scheme.regenerator.HotInput().PressureInput(), scheme.regenerator.HotInput().MassRateInput(),
		},
	)
	sink.SinkAll(
		// ports on the first line are sinked cos parametric turbine does not take output pressure
		// as an input parameter
		scheme.regenerator.HotOutput().GasOutput(),
		scheme.regenerator.HotOutput().TemperatureOutput(),
		scheme.regenerator.HotOutput().PressureOutput(),
		scheme.regenerator.HotOutput().MassRateOutput(),
	)

	graph.Link(
		scheme.payload.RPMOutput(), scheme.fTurbine.RPMInput(),
	)

	sink.SinkAll(scheme.payload.PowerOutput(), scheme.fTurbine.PowerOutput())
}

func (scheme *doubleShaftRegFreeScheme) setEquations() {
	scheme.gasGenMassRateEq = utils.NewEquality(
		graph.NewWeakPort(scheme.gasGeneratorPart.Burner.MassRateOutput()),
		graph.NewWeakPort(scheme.gasGeneratorPart.Turbine.MassRateInput()),
	)
	scheme.gasGenMassRateEq.SetName("gasGenMassRateEq")

	scheme.gasGenPowerEq = utils.NewMultiAdderFromPorts(
		[]graph.Port{
			graph.NewWeakPort(scheme.gasGeneratorPart.Turbine.PowerOutput()),
			graph.NewWeakPort(scheme.gasGeneratorPart.Turbine.MassRateInput()),
		},
		[]graph.Port{
			graph.NewWeakPort(scheme.gasGeneratorPart.Shaft.PowerOutput()),
			graph.NewWeakPort(scheme.gasGeneratorPart.Compressor.MassRateInput()),
		},
	)
	scheme.gasGenPowerEq.SetName("gasGenPowerEq")

	scheme.freeTurbineMassRateEq = utils.NewEquality(
		graph.NewWeakPort(scheme.ctPipe.MassRateOutput()),
		graph.NewWeakPort(scheme.fTurbine.MassRateInput()),
	)
	scheme.freeTurbineMassRateEq.SetName("freeTurbineMassRateEq")

	scheme.freeTurbinePowerEq = utils.NewMultiAdderFromPorts(
		[]graph.Port{
			graph.NewWeakPort(scheme.fTurbine.MassRateInput()),
			graph.NewWeakPort(scheme.fTurbine.PowerOutput()),
		},
		[]graph.Port{graph.NewWeakPort(scheme.payload.PowerOutput())},
	)
	scheme.freeTurbinePowerEq.SetName("freeTurbinePowerEq")

	scheme.freeTurbinePressureEq = utils.NewEquality(
		graph.NewWeakPort(scheme.fTurbine.PressureOutput()),
		graph.NewWeakPort(scheme.ftPipe.PressureInput()),
	)
	scheme.freeTurbinePressureEq.SetName("freeTurbinePressureEq")

	scheme.gasGenBurnerEq = utils.NewEquality(
		scheme.burnerTemperatureSource.TemperatureOutput(),
		graph.NewWeakPort(scheme.gasGeneratorPart.Burner.TemperatureOutput()),
	)
	scheme.gasGenBurnerEq.SetName("gasGenBurnerEq")

	scheme.assembler.AddInputPorts(
		scheme.gasGenMassRateEq.OutputPort(),
		scheme.gasGenPowerEq.OutputPort(),
		scheme.freeTurbineMassRateEq.OutputPort(),
		scheme.freeTurbinePowerEq.OutputPort(),
		scheme.freeTurbinePressureEq.OutputPort(),
		scheme.gasGenBurnerEq.OutputPort(),
	)
	sink.SinkPort(scheme.assembler.GetVectorPort())
}
