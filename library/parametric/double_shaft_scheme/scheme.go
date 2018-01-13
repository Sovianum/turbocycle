package parametric

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/variator"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/material/gases"
)

type DoubleShaftFreeScheme interface {
	Compressor() c.ParametricCompressorNode
	CompressorPipe() c.PressureLossNode
	Burner() c.ParametricBurnerNode
	CompressorTurbine() c.ParametricTurbineNode
	CompressorTurbinePipe() c.PressureLossNode
	FreeTurbine() c.ParametricTurbineNode
	FreeTurbinePipe() c.PressureLossNode
	Payload() c.Payload
	Assembler() graph.VectorAssemblerNode
	Variators() []variator.Variator
	GetNetwork() (graph.Network, error)
}

func NewDoubleShaftFreeScheme(
	gas gases.Gas, tAtm, pAtm, tGas, etaM float64,
	compressor c.ParametricCompressorNode,
	compressorPipe c.PressureLossNode,
	burner c.ParametricBurnerNode,
	compressorTurbine c.ParametricTurbineNode,
	shaft c.TransmissionNode,
	compressorTurbinePipe c.PressureLossNode,
	freeTurbine c.ParametricTurbineNode,
	freeTurbinePipe c.PressureLossNode,
	payload c.Payload,
) DoubleShaftFreeScheme {
	var result = &doubleShaftFreeScheme{
		gasPart: newGasPart(gas, tAtm, pAtm, pAtm),
		gasGeneratorPart: newGasGeneratorPart(
			compressor, burner, compressorTurbine, shaft, compressorPipe,
		),

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

type doubleShaftFreeScheme struct {
	gasPart          *gasPart
	gasGeneratorPart *gasGeneratorPart

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

func (scheme *doubleShaftFreeScheme) Compressor() c.ParametricCompressorNode {
	return scheme.gasGeneratorPart.compressor
}

func (scheme *doubleShaftFreeScheme) CompressorPipe() c.PressureLossNode {
	return scheme.gasGeneratorPart.compressorPipe
}

func (scheme *doubleShaftFreeScheme) Burner() c.ParametricBurnerNode {
	return scheme.gasGeneratorPart.burner
}

func (scheme *doubleShaftFreeScheme) CompressorTurbine() c.ParametricTurbineNode {
	return scheme.gasGeneratorPart.turbine
}

func (scheme *doubleShaftFreeScheme) CompressorTurbinePipe() c.PressureLossNode {
	return scheme.ctPipe
}

func (scheme *doubleShaftFreeScheme) FreeTurbine() c.ParametricTurbineNode {
	return scheme.fTurbine
}

func (scheme *doubleShaftFreeScheme) FreeTurbinePipe() c.PressureLossNode {
	return scheme.ftPipe
}

func (scheme *doubleShaftFreeScheme) Payload() c.Payload {
	return scheme.payload
}

func (scheme *doubleShaftFreeScheme) Assembler() graph.VectorAssemblerNode {
	return scheme.assembler
}

func (scheme *doubleShaftFreeScheme) Variators() []variator.Variator {
	return scheme.variators
}

func (scheme *doubleShaftFreeScheme) GetNetwork() (graph.Network, error) {
	var nodes = append(scheme.gasPart.Nodes(), scheme.gasGeneratorPart.Nodes()...)
	nodes = append(
		nodes, scheme.ctPipe, scheme.ftPipe, scheme.fTurbine, scheme.payload, scheme.burnerTemperatureSource,
		scheme.assembler,
	)
	return graph.NewNetwork(nodes)
}

func (scheme *doubleShaftFreeScheme) linkPorts() {
	graph.LinkAll(
		[]graph.Port{
			scheme.gasPart.gasSource.GasOutput(),
			scheme.gasPart.temperatureSource.TemperatureOutput(),
			scheme.gasPart.inputPressureSource.PressureOutput(),
		},
		[]graph.Port{
			scheme.gasGeneratorPart.compressor.GasInput(),
			scheme.gasGeneratorPart.compressor.TemperatureInput(),
			scheme.gasGeneratorPart.compressor.PressureInput(),
		},
	)
	sink.SinkPort(scheme.gasGeneratorPart.compressor.MassRateInput())

	graph.LinkAll(
		[]graph.Port{
			scheme.gasGeneratorPart.turbine.GasOutput(),
			scheme.gasGeneratorPart.turbine.TemperatureOutput(),
			scheme.gasGeneratorPart.turbine.PressureOutput(),
			scheme.gasGeneratorPart.turbine.MassRateOutput(),
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
	sink.SinkAll(
		// ports on the first line are sinked cos parametric turbine does not take output pressure
		// as an input parameter
		scheme.fTurbine.PressureOutput(), scheme.ftPipe.PressureInput(),
		scheme.ftPipe.GasOutput(), scheme.ftPipe.TemperatureOutput(),
		scheme.ftPipe.MassRateOutput(),
	)

	graph.Link(
		scheme.gasPart.outputPressureSource.PressureOutput(),
		scheme.ftPipe.PressureOutput(),
	)

	graph.Link(
		scheme.payload.RPMOutput(), scheme.fTurbine.RPMInput(),
	)

	sink.SinkAll(scheme.payload.PowerOutput())
}

func (scheme *doubleShaftFreeScheme) setEquations() {
	scheme.gasGenMassRateEq = c.NewEquality(
		graph.NewWeakPort(scheme.gasGeneratorPart.burner.MassRateOutput()),
		graph.NewWeakPort(scheme.gasGeneratorPart.turbine.MassRateInput()),
	)
	scheme.gasGenMassRateEq.SetName("gasGenMassRateEq")

	scheme.gasGenPowerEq = c.NewMultiAdderFromPorts(
		[]graph.Port{
			graph.NewWeakPort(scheme.gasGeneratorPart.turbine.PowerOutput()),
			graph.NewWeakPort(scheme.gasGeneratorPart.turbine.MassRateOutput()),
		},
		[]graph.Port{
			graph.NewWeakPort(scheme.gasGeneratorPart.compressor.PowerOutput()),
			graph.NewWeakPort(scheme.gasGeneratorPart.compressor.MassRateOutput()),
		},
	)
	scheme.gasGenPowerEq.SetName("gasGenPowerEq")

	scheme.freeTurbineMassRateEq = c.NewEquality(
		graph.NewWeakPort(scheme.ctPipe.MassRateOutput()),
		graph.NewWeakPort(scheme.fTurbine.MassRateInput()),
	)
	scheme.freeTurbineMassRateEq.SetName("freeTurbineMassRateEq")

	scheme.freeTurbinePowerEq = c.NewMultiAdderFromPorts(
		[]graph.Port{
			graph.NewWeakPort(scheme.fTurbine.MassRateOutput()),
			graph.NewWeakPort(scheme.fTurbine.PowerOutput()),
		},
		[]graph.Port{graph.NewWeakPort(scheme.payload.PowerOutput())},
	)
	scheme.freeTurbinePowerEq.SetName("freeTurbinePowerEq")

	scheme.freeTurbinePressureEq = c.NewEquality(
		graph.NewWeakPort(scheme.fTurbine.PressureOutput()),
		graph.NewWeakPort(scheme.ftPipe.PressureInput()),
	)
	scheme.freeTurbinePressureEq.SetName("freeTurbinePressureEq")

	scheme.gasGenBurnerEq = c.NewEquality(
		scheme.burnerTemperatureSource.TemperatureOutput(),
		graph.NewWeakPort(scheme.gasGeneratorPart.burner.TemperatureOutput()),
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
