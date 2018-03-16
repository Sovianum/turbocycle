package free3n

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/variator"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive/utils"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/library/parametric"
	"github.com/Sovianum/turbocycle/material/gases"
)

type ThreeShaftFreeScheme interface {
	TemperatureSource() source.TemperatureSourceNode

	HPC() c.ParametricCompressorNode
	HPCPipe() c.PressureLossNode
	HPT() c.ParametricTurbineNode
	HPTPipe() c.PressureLossNode
	HPShaft() c.TransmissionNode

	LPC() c.ParametricCompressorNode
	LPCPipe() c.PressureLossNode
	LPT() c.ParametricTurbineNode
	LPTPipe() c.PressureLossNode
	LPShaft() c.TransmissionNode

	Burner() c.ParametricBurnerNode

	FT() c.ParametricTurbineNode
	FTPipe() c.PressureLossNode

	Payload() c.Payload
	Assembler() graph.VectorAssemblerNode
	Variators() []variator.Variator
	GetNetwork() (graph.Network, error)
}

func NewThreeShaftFreeScheme(
	gas gases.Gas, tAtm, pAtm, tGas float64,
	mpc c.ParametricCompressorNode, mpcPipe c.PressureLossNode,
	mpt c.ParametricTurbineNode, mptPipe c.PressureLossNode,
	mpEtaM float64,

	hpc c.ParametricCompressorNode, hpcPipe c.PressureLossNode,
	hpt c.ParametricTurbineNode, hptPipe c.PressureLossNode,
	hpEtaM float64,

	lpt c.ParametricTurbineNode, lptPipe c.PressureLossNode,
	burner c.ParametricBurnerNode, payload c.Payload,
) ThreeShaftFreeScheme {
	mpc.SetName("mpc")
	mpcPipe.SetName("mpcPipe")
	mpt.SetName("mpt")
	mptPipe.SetName("mptPipe")

	hpc.SetName("hpc")
	hpcPipe.SetName("hpcPipe")
	hpt.SetName("hpt")
	hptPipe.SetName("hptPipe")

	lpt.SetName("lpt")
	lptPipe.SetName("lptPipe")

	result := &threeShaftFreeScheme{
		gasPart: parametric.NewGasPart(gas, tAtm, pAtm, pAtm),

		burner: burner,

		hpShaft: parametric.NewTurboShaftPart(hpc, hpt, c.NewTransmissionNode(hpEtaM)),
		hpcPipe: hpcPipe,
		hptPipe: hptPipe,

		mpShaft: parametric.NewTurboShaftPart(mpc, mpt, c.NewTransmissionNode(mpEtaM)),
		mpcPipe: mpcPipe,
		mptPipe: mptPipe,

		lpt:     lpt,
		lptPipe: lptPipe,

		payload: payload,

		burnerTSource: source.NewTemperatureSourceNode(tGas),

		assembler: graph.NewVectorAssemblerNode(),
		variators: []variator.Variator{
			variator.FromCallables(mpc.NormMassRate, mpc.SetNormMassRate),
			variator.FromCallables(mpc.NormPiStag, mpc.SetNormPiStag),

			variator.FromCallables(hpc.NormMassRate, hpc.SetNormMassRate),
			variator.FromCallables(hpc.NormPiStag, hpc.SetNormPiStag),

			variator.FromCallables(lpt.NormPiT, lpt.SetNormPiT),
			variator.FromCallables(mpt.NormPiT, mpt.SetNormPiT),
			variator.FromCallables(hpt.NormPiT, hpt.SetNormPiT),

			variator.FromCallables(burner.FuelRateRel, burner.SetFuelRateRel),
			variator.FromCallables(payload.NormRPM, payload.SetNormRPM),
		},
	}
	result.linkPorts()
	result.setEquations()
	return result
}

type threeShaftFreeScheme struct {
	gasPart *parametric.GasPart

	burner c.ParametricBurnerNode

	hpShaft *parametric.TurboShaftPart
	hpcPipe c.PressureLossNode
	hptPipe c.PressureLossNode

	mpShaft *parametric.TurboShaftPart
	mpcPipe c.PressureLossNode
	mptPipe c.PressureLossNode

	lpt     c.ParametricTurbineNode
	lptPipe c.PressureLossNode

	payload c.Payload

	burnerTSource source.TemperatureSourceNode

	compMassRateEq graph.ReduceNode
	hptMassRateEq  graph.ReduceNode
	mptMassRateEq  graph.ReduceNode
	lptMassRateEq  graph.ReduceNode

	hpPowerEq graph.ReduceNode
	mpPowerEq graph.ReduceNode
	lpPowerEq graph.ReduceNode

	lpPressureEq graph.ReduceNode
	burnerEq     graph.ReduceNode

	assembler graph.VectorAssemblerNode
	variators []variator.Variator
}

func (scheme *threeShaftFreeScheme) HPShaft() c.TransmissionNode {
	return scheme.hpShaft.Shaft
}

func (scheme *threeShaftFreeScheme) LPShaft() c.TransmissionNode {
	return scheme.mpShaft.Shaft
}

func (scheme *threeShaftFreeScheme) TemperatureSource() source.TemperatureSourceNode {
	return scheme.burnerTSource
}

func (scheme *threeShaftFreeScheme) HPC() c.ParametricCompressorNode {
	return scheme.hpShaft.Compressor
}

func (scheme *threeShaftFreeScheme) HPCPipe() c.PressureLossNode {
	return scheme.hpcPipe
}

func (scheme *threeShaftFreeScheme) Burner() c.ParametricBurnerNode {
	return scheme.burner
}

func (scheme *threeShaftFreeScheme) HPT() c.ParametricTurbineNode {
	return scheme.hpShaft.Turbine
}

func (scheme *threeShaftFreeScheme) HPTPipe() c.PressureLossNode {
	return scheme.hptPipe
}

func (scheme *threeShaftFreeScheme) LPC() c.ParametricCompressorNode {
	return scheme.mpShaft.Compressor
}

func (scheme *threeShaftFreeScheme) LPCPipe() c.PressureLossNode {
	return scheme.mpcPipe
}

func (scheme *threeShaftFreeScheme) LPT() c.ParametricTurbineNode {
	return scheme.mpShaft.Turbine
}

func (scheme *threeShaftFreeScheme) LPTPipe() c.PressureLossNode {
	return scheme.mptPipe
}

func (scheme *threeShaftFreeScheme) FT() c.ParametricTurbineNode {
	return scheme.lpt
}

func (scheme *threeShaftFreeScheme) FTPipe() c.PressureLossNode {
	return scheme.lptPipe
}

func (scheme *threeShaftFreeScheme) Payload() c.Payload {
	return scheme.payload
}

func (scheme *threeShaftFreeScheme) Assembler() graph.VectorAssemblerNode {
	return scheme.assembler
}

func (scheme *threeShaftFreeScheme) Variators() []variator.Variator {
	return scheme.variators
}

func (scheme *threeShaftFreeScheme) GetNetwork() (graph.Network, error) {
	return graph.NewNetwork(scheme.getNodes())
}

func (scheme *threeShaftFreeScheme) getNodes() []graph.Node {
	var nodes = append(
		scheme.gasPart.Nodes(),
		scheme.payload, scheme.burnerTSource,
		scheme.lpt, scheme.lptPipe,
		scheme.mpcPipe, scheme.mptPipe,
		scheme.hpcPipe, scheme.hptPipe,
		scheme.burner,
	)
	nodes = append(nodes, scheme.hpShaft.Nodes()...)
	nodes = append(nodes, scheme.mpShaft.Nodes()...)
	nodes = append(
		nodes,
		scheme.compMassRateEq, scheme.hptMassRateEq,
		scheme.mptMassRateEq, scheme.lptMassRateEq,

		scheme.hpPowerEq, scheme.mpPowerEq, scheme.lpPowerEq,
		scheme.lpPressureEq, scheme.burnerEq,
		scheme.assembler,
	)
	return nodes
}

func (scheme *threeShaftFreeScheme) linkPorts() {
	scheme.linkHPShaft()
	scheme.linkMPShaft()
	scheme.linkLPShaft()
}

func (scheme *threeShaftFreeScheme) linkHPShaft() {
	graph.LinkAll(
		[]graph.Port{
			scheme.mpcPipe.GasOutput(),
			scheme.mpcPipe.TemperatureOutput(),
			scheme.mpcPipe.PressureOutput(),
		},
		[]graph.Port{
			scheme.hpShaft.Compressor.GasInput(),
			scheme.hpShaft.Compressor.TemperatureInput(),
			scheme.hpShaft.Compressor.PressureInput(),
		},
	)
	sink.SinkAll(
		scheme.mpcPipe.MassRateOutput(),
		scheme.hpShaft.Compressor.MassRateInput(),
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.hpShaft.Compressor.GasOutput(),
			scheme.hpShaft.Compressor.TemperatureOutput(),
			scheme.hpShaft.Compressor.PressureOutput(),
			scheme.hpShaft.Compressor.MassRateOutput(),
		},
		[]graph.Port{
			scheme.hpcPipe.GasInput(),
			scheme.hpcPipe.TemperatureInput(),
			scheme.hpcPipe.PressureInput(),
			scheme.hpcPipe.MassRateInput(),
		},
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.hpcPipe.GasOutput(),
			scheme.hpcPipe.TemperatureOutput(),
			scheme.hpcPipe.PressureOutput(),
			scheme.hpcPipe.MassRateOutput(),
		},
		[]graph.Port{
			scheme.burner.GasInput(),
			scheme.burner.TemperatureInput(),
			scheme.burner.PressureInput(),
			scheme.burner.MassRateInput(),
		},
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.burner.GasOutput(),
			scheme.burner.TemperatureOutput(),
			scheme.burner.PressureOutput(),
		},
		[]graph.Port{
			scheme.hpShaft.Turbine.GasInput(),
			scheme.hpShaft.Turbine.TemperatureInput(),
			scheme.hpShaft.Turbine.PressureInput(),
		},
	)
	sink.SinkAll(
		scheme.burner.MassRateOutput(), scheme.hpShaft.Turbine.MassRateInput(),
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.hpShaft.Turbine.GasOutput(),
			scheme.hpShaft.Turbine.TemperatureOutput(),
			scheme.hpShaft.Turbine.PressureOutput(),
			scheme.hpShaft.Turbine.MassRateOutput(),
		},
		[]graph.Port{
			scheme.hptPipe.GasInput(),
			scheme.hptPipe.TemperatureInput(),
			scheme.hptPipe.PressureInput(),
			scheme.hptPipe.MassRateInput(),
		},
	)
}

func (scheme *threeShaftFreeScheme) linkMPShaft() {
	graph.LinkAll(
		[]graph.Port{
			scheme.gasPart.GasSource.GasOutput(),
			scheme.gasPart.TemperatureSource.TemperatureOutput(),
			scheme.gasPart.InputPressureSource.PressureOutput(),
		},
		[]graph.Port{
			scheme.mpShaft.Compressor.GasInput(),
			scheme.mpShaft.Compressor.TemperatureInput(),
			scheme.mpShaft.Compressor.PressureInput(),
		},
	)
	sink.SinkAll(
		scheme.mpShaft.Compressor.MassRateInput(),
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.mpShaft.Compressor.GasOutput(),
			scheme.mpShaft.Compressor.TemperatureOutput(),
			scheme.mpShaft.Compressor.PressureOutput(),
			scheme.mpShaft.Compressor.MassRateOutput(),
		},
		[]graph.Port{
			scheme.mpcPipe.GasInput(),
			scheme.mpcPipe.TemperatureInput(),
			scheme.mpcPipe.PressureInput(),
			scheme.mpcPipe.MassRateInput(),
		},
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.hptPipe.GasOutput(),
			scheme.hptPipe.TemperatureOutput(),
			scheme.hptPipe.PressureOutput(),
		},
		[]graph.Port{
			scheme.mpShaft.Turbine.GasInput(),
			scheme.mpShaft.Turbine.TemperatureInput(),
			scheme.mpShaft.Turbine.PressureInput(),
		},
	)
	sink.SinkAll(
		scheme.hptPipe.MassRateOutput(), scheme.mpShaft.Turbine.MassRateInput(),
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.mpShaft.Turbine.GasOutput(),
			scheme.mpShaft.Turbine.TemperatureOutput(),
			scheme.mpShaft.Turbine.PressureOutput(),
			scheme.mpShaft.Turbine.MassRateOutput(),
		},
		[]graph.Port{
			scheme.mptPipe.GasInput(),
			scheme.mptPipe.TemperatureInput(),
			scheme.mptPipe.PressureInput(),
			scheme.mptPipe.MassRateInput(),
		},
	)
}

func (scheme *threeShaftFreeScheme) linkLPShaft() {
	graph.LinkAll(
		[]graph.Port{
			scheme.mptPipe.GasOutput(),
			scheme.mptPipe.TemperatureOutput(),
			scheme.mptPipe.PressureOutput(),
		},
		[]graph.Port{
			scheme.lpt.GasInput(),
			scheme.lpt.TemperatureInput(),
			scheme.lpt.PressureInput(),
		},
	)
	graph.Link(
		scheme.gasPart.OutputPressureSource.PressureOutput(),
		scheme.lptPipe.PressureOutput(),
	)
	sink.SinkAll(
		scheme.mptPipe.MassRateOutput(), scheme.lpt.MassRateInput(),
		scheme.lpt.GasOutput(), scheme.lpt.TemperatureOutput(),
		scheme.lpt.PressureOutput(), scheme.lpt.MassRateOutput(),
		scheme.lpt.PowerOutput(),
	)
	sink.SinkAll(
		scheme.lptPipe.GasOutput(), scheme.lptPipe.TemperatureOutput(),
		scheme.lptPipe.MassRateOutput(),
		scheme.lptPipe.GasInput(), scheme.lptPipe.TemperatureInput(),
		scheme.lptPipe.PressureInput(), scheme.lptPipe.MassRateInput(),
		scheme.payload.PowerOutput(),
	)
	graph.Link(scheme.payload.RPMOutput(), scheme.lpt.RPMInput())
}

func (scheme *threeShaftFreeScheme) setEquations() {
	scheme.compMassRateEq = utils.NewEquality(
		graph.NewWeakPort(scheme.mpcPipe.MassRateOutput()),
		graph.NewWeakPort(scheme.hpShaft.Compressor.MassRateInput()),
	)
	scheme.compMassRateEq.SetName("compMassRateEq")

	scheme.hptMassRateEq = utils.NewEquality(
		graph.NewWeakPort(scheme.burner.MassRateOutput()),
		graph.NewWeakPort(scheme.hpShaft.Turbine.MassRateInput()),
	)
	scheme.hptMassRateEq.SetName("hptMassRateEq")

	scheme.mptMassRateEq = utils.NewEquality(
		graph.NewWeakPort(scheme.hptPipe.MassRateOutput()),
		graph.NewWeakPort(scheme.mpShaft.Turbine.MassRateInput()),
	)
	scheme.mptMassRateEq.SetName("mptMassRateEq")

	scheme.lptMassRateEq = utils.NewEquality(
		graph.NewWeakPort(scheme.mptPipe.MassRateOutput()),
		graph.NewWeakPort(scheme.lpt.MassRateInput()),
	)
	scheme.lptMassRateEq.SetName("lptMassRateEq")

	scheme.hpPowerEq = utils.NewMultiAdderFromPorts(
		[]graph.Port{
			graph.NewWeakPort(scheme.hpShaft.Compressor.MassRateInput()),
			graph.NewWeakPort(scheme.hpShaft.Shaft.PowerOutput()),
		},
		[]graph.Port{
			graph.NewWeakPort(scheme.hpShaft.Turbine.MassRateInput()),
			graph.NewWeakPort(scheme.hpShaft.Turbine.PowerOutput()),
		},
	)
	scheme.hpPowerEq.SetName("hpPowerEq")

	scheme.mpPowerEq = utils.NewMultiAdderFromPorts(
		[]graph.Port{
			graph.NewWeakPort(scheme.mpShaft.Compressor.MassRateInput()),
			graph.NewWeakPort(scheme.mpShaft.Shaft.PowerOutput()),
		},
		[]graph.Port{
			graph.NewWeakPort(scheme.mpShaft.Turbine.MassRateInput()),
			graph.NewWeakPort(scheme.mpShaft.Turbine.PowerOutput()),
		},
	)
	scheme.mpPowerEq.SetName("mpPowerEq")

	scheme.lpPowerEq = utils.NewMultiAdderFromPorts(
		[]graph.Port{
			graph.NewWeakPort(scheme.lpt.MassRateInput()),
			graph.NewWeakPort(scheme.lpt.PowerOutput()),
		},
		[]graph.Port{
			graph.NewWeakPort(scheme.payload.PowerOutput()),
		},
	)
	scheme.lpPowerEq.SetName("lpPowerEq")

	scheme.lpPressureEq = utils.NewEquality(
		graph.NewWeakPort(scheme.lptPipe.PressureInput()),
		graph.NewWeakPort(scheme.lpt.PressureOutput()),
	)
	scheme.lpPressureEq.SetName("lpPressureEq")

	scheme.burnerEq = utils.NewEquality(
		scheme.burnerTSource.TemperatureOutput(),
		graph.NewWeakPort(scheme.burner.TemperatureOutput()),
	)
	scheme.burnerEq.SetName("burnerEq")

	scheme.assembler.AddInputPorts(
		scheme.compMassRateEq.OutputPort(),
		scheme.hptMassRateEq.OutputPort(),
		scheme.mptMassRateEq.OutputPort(),
		scheme.lptMassRateEq.OutputPort(),

		scheme.hpPowerEq.OutputPort(),
		scheme.mpPowerEq.OutputPort(),
		scheme.lpPowerEq.OutputPort(),

		scheme.lpPressureEq.OutputPort(),
		scheme.burnerEq.OutputPort(),
	)
	sink.SinkPort(scheme.assembler.GetVectorPort())
}
