package parametric

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/core/math/variator"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/material/gases"
)

type SingleShaftScheme interface {
	Compressor() c.ParametricCompressorNode
	Burner() c.ParametricBurnerNode
	Turbine() c.ParametricTurbineNode
	Payload() c.Payload
	GetNetwork() (graph.Network, error)
}

func NewSingleShaftScheme(
	gas gases.Gas, pAtm, tAtm, tGas float64,
	compressor c.ParametricCompressorNode,
	burner c.ParametricBurnerNode,
	turbine c.ParametricTurbineNode,
	payload c.Payload,
) SingleShaftScheme {
	var result = &singleShaftScheme{
		compressor: compressor,
		burner:     burner,
		turbine:    turbine,
		payload:    payload,

		hub: helper.NewHubNode(),

		gasSource:            source.NewGasSourceNode(gas),
		temperatureSource:    source.NewTemperatureSourceNode(tAtm),
		inputPressureSource:  source.NewPressureSourceNode(pAtm),
		outputPressureSource: source.NewPressureSourceNode(pAtm),

		burnerTemperatureSource: source.NewTemperatureSourceNode(tGas),

		assembler: graph.NewVectorAssemblerNode(),

		variators: []variator.Variator{
			variator.FromCallables(compressor.NormMassRate, compressor.SetNormMassRate),
			variator.FromCallables(compressor.NormPiStag, compressor.SetNormPiStag),
			variator.FromCallables(burner.FuelRateRel, burner.SetFuelRateRel),
			variator.FromCallables(payload.NormRPM, payload.SetNormRPM),
			variator.FromCallables(turbine.NormPiT, turbine.SetNormPiT),
		},
	}
	result.linkPorts()
	result.setEquations()

	return result
}

type singleShaftScheme struct {
	compressor c.ParametricCompressorNode
	burner     c.ParametricBurnerNode
	turbine    c.ParametricTurbineNode
	payload    c.Payload

	hub helper.HubNode

	gasSource            source.GasSourceNode
	temperatureSource    source.TemperatureSourceNode
	inputPressureSource  source.PressureSourceNode
	outputPressureSource source.PressureSourceNode

	burnerTemperatureSource source.TemperatureSourceNode

	massRateEq    graph.ReduceNode
	powerEq       graph.ReduceNode
	rpmEq         graph.ReduceNode
	pressureEq    graph.ReduceNode
	temperatureEq graph.ReduceNode

	assembler graph.VectorAssemblerNode

	variators []variator.Variator

	solver math.Solver
}

func (scheme *singleShaftScheme) GetNetwork() (graph.Network, error) {
	var nodes = []graph.Node{
		scheme.compressor, scheme.burner, scheme.turbine, scheme.payload,
		scheme.hub, scheme.gasSource, scheme.temperatureSource,
		scheme.inputPressureSource, scheme.outputPressureSource,
		scheme.burnerTemperatureSource, scheme.massRateEq, scheme.powerEq,
		scheme.rpmEq, scheme.pressureEq, scheme.temperatureEq,
		scheme.assembler,
	}

	return graph.NewNetwork(nodes)
}

func (scheme *singleShaftScheme) Compressor() c.ParametricCompressorNode {
	return scheme.compressor
}

func (scheme *singleShaftScheme) Burner() c.ParametricBurnerNode {
	return scheme.burner
}

func (scheme *singleShaftScheme) Turbine() c.ParametricTurbineNode {
	return scheme.turbine
}

func (scheme *singleShaftScheme) Payload() c.Payload {
	return scheme.payload
}

func (scheme *singleShaftScheme) linkPorts() {
	graph.LinkAll(
		[]graph.Port{
			scheme.inputPressureSource.PressureOutput(),
			scheme.temperatureSource.TemperatureOutput(),
			scheme.gasSource.GasOutput(),
		},
		[]graph.Port{
			scheme.compressor.GasInput(),
			scheme.compressor.TemperatureInput(),
			scheme.compressor.PressureInput(),
		},
	)

	graph.LinkAll(
		[]graph.Port{
			scheme.compressor.GasOutput(),
			scheme.compressor.TemperatureOutput(),
			scheme.compressor.PressureOutput(),
			scheme.compressor.MassRateOutput(),
		},
		[]graph.Port{
			scheme.burner.GasInput(),
			scheme.burner.TemperatureInput(),
			scheme.burner.PressureInput(),
			scheme.burner.MassRateInput(),
		},
	)

	graph.Link(scheme.burner.TemperatureOutput(), scheme.hub.Inlet())

	graph.LinkAll(
		[]graph.Port{
			scheme.burner.GasOutput(),
			scheme.hub.Outlet1(),
			scheme.burner.PressureOutput(),
		},
		[]graph.Port{
			scheme.turbine.GasInput(),
			scheme.turbine.TemperatureInput(),
			scheme.turbine.PressureInput(),
		},
	)

	graph.Link(scheme.turbine.RPMInput(), scheme.compressor.RPMOutput())

	sink.SinkPort(scheme.compressor.MassRateInput())
	sink.SinkPort(scheme.turbine.TemperatureOutput())
	sink.SinkPort(scheme.turbine.GasOutput())
}

func (scheme *singleShaftScheme) setEquations() {
	scheme.massRateEq = c.NewEquality(
		scheme.burner.MassRateOutput(),
		scheme.turbine.MassRateInput(),
	)
	// todo check sign of specific labour of payload
	scheme.powerEq = c.NewMultiAdderFromPorts([][]graph.Port{
		{scheme.turbine.PowerOutput(), scheme.turbine.MassRateOutput()},
		{scheme.compressor.PowerOutput(), scheme.compressor.MassRateOutput()},
		{scheme.payload.PowerOutput()},
	})
	scheme.rpmEq = c.NewEquality(
		scheme.payload.RPMOutput(),
		scheme.turbine.RPMInput(),
	)
	scheme.pressureEq = c.NewEquality(
		scheme.turbine.PressureOutput(),
		scheme.outputPressureSource.PressureOutput(),
	)
	scheme.temperatureEq = c.NewEquality(
		scheme.burnerTemperatureSource.TemperatureOutput(),
		scheme.hub.Outlet2(),
	)

	scheme.assembler.AddInputPorts(
		scheme.massRateEq.OutputPort(),
		scheme.powerEq.OutputPort(),
		scheme.rpmEq.OutputPort(),
		scheme.pressureEq.OutputPort(),
		scheme.temperatureEq.OutputPort(),
	)
	sink.SinkPort(scheme.assembler.GetVectorPort())
}
