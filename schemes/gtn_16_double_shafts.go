package schemes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/nodes/source"
	"github.com/Sovianum/turbocycle/impl/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/impl/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/states"
)

const (
	tFuel = 300
	t0 = 288
)

func NewGtn16DoubleShaft(
	atmGas gases.Gas, pAtm, tAtm float64,
	etaCompressor, piStag float64,
	etaTBlocked, lambdaOutBlocked, dMassRateRelBlocked float64,
	fuel fuel.GasFuel, tGasStag, sigmaBurn, etaBurn, initAlpha float64,
	etaTFree, lambdaOutFree, dMassRateRelFree float64,

	nodePrecision float64,
) Gtn16DoubleShaft {
	return &gtn16DoubleShaft{
		atmGas:atmGas,
		pAtm:pAtm,
		tAtm:tAtm,

		etaCompressor:etaCompressor,
		piStag:piStag,

		etaTBlocked:etaTBlocked,
		lambdaOutBlocked:lambdaOutBlocked,
		dMassRateRelBlocked:dMassRateRelBlocked,

		fuel:fuel,
		tGasStag:tGasStag,
		sigmaBurn:sigmaBurn,
		etaBurn:etaBurn,
		initAlpha:initAlpha,

		etaTFree:etaTFree,
		lambdaOutFree:lambdaOutFree,
		dMassRateRelFree:dMassRateRelFree,

		nodePrecision:nodePrecision,
	}
}

type Gtn16DoubleShaft interface {
	Scheme
}

type gtn16DoubleShaft struct {
	atmGas gases.Gas
	tAtm float64
	pAtm float64

	etaCompressor float64
	piStag float64

	etaTBlocked float64
	lambdaOutBlocked float64
	dMassRateRelBlocked float64

	fuel fuel.GasFuel
	tGasStag float64
	sigmaBurn float64
	etaBurn float64
	initAlpha float64

	etaTFree float64
	lambdaOutFree float64
	dMassRateRelFree float64

	nodePrecision float64
}

//def _link_nodes(self):
//# connect gas channel
//self.link_ports(self.atm_node.get_new_output_port(), self.inlet_filter_node.gas_input)
//self.link_ports(self.inlet_filter_node.gas_output, self.inlet_pipe_node.gas_input)
//self.link_ports(self.inlet_pipe_node.gas_output, self.compressor_node.gas_input)
//self.link_ports(self.compressor_node.gas_output, self.regenerator_node.cold_input)
//self.link_ports(self.regenerator_node.cold_output, self.burner_node.gas_input)
//self.link_ports(self.burner_node.gas_output, self.compressor_turbine_node.gas_input)
//self.link_ports(self.compressor_turbine_node.gas_output, self.compressor_turbine_pipe_node.gas_input)
//self.link_ports(self.compressor_turbine_pipe_node.gas_output, self.free_turbine_node.gas_input)
//self.link_ports(self.free_turbine_node.gas_output, self.free_turbine_tube_node.gas_input)
//self.link_ports(self.free_turbine_tube_node.gas_output, self.regenerator_node.hot_input)
//self.link_ports(self.regenerator_node.hot_output, self.regenerator_tube_node.gas_input)
//self.link_ports(self.regenerator_tube_node.gas_output, self.atm_node.get_new_dynamic_input_port())
//
//# connect power channel
//self.link_ports(self.compressor_node.power_output, self.transmission_node.power_input)
//self.link_ports(self.transmission_node.power_output, self.compressor_turbine_node.power_input)
//
//# collect all power outputs by collector
//self.link_ports(self.compressor_node.to_collector, self.collector_node.get_new_input_port())
//self.link_ports(self.compressor_turbine_node.power_output, self.collector_node.get_new_input_port())
//self.link_ports(self.free_turbine_node.power_output, self.collector_node.get_new_input_port())

func (scheme *gtn16DoubleShaft) GetNetwork() *core.Network {
	var nodes = make(map[string]core.Node)

	var inputGasSource = source.NewComplexGasSourceNode(scheme.atmGas, scheme.tAtm, scheme.pAtm)
	nodes[inputGasSourceName] = inputGasSource
	var outputGasSource = source.NewComplexGasSourceNode(scheme.atmGas, scheme.tAtm, scheme.pAtm)
	nodes[outputGasSourceName] = outputGasSource
	var gasSink = sink.NewComplexGasSinkNode()
	nodes[outputGasSinkName] = gasSink
	var compressor = constructive.NewCompressorNode(scheme.etaCompressor, scheme.piStag, scheme.nodePrecision)
	nodes[compressorName] = compressor
	var compressorTurbine = constructive.NewBlockedTurbineNode(scheme.etaTBlocked, scheme.lambdaOutBlocked, scheme.nodePrecision,
		func(constructive.TurbineNode) float64 {
			return scheme.dMassRateRelBlocked
		},
	)
	nodes[compressorTurbineName] = compressorTurbine
	var burner = constructive.NewBurnerNode(
		scheme.fuel, scheme.tGasStag, tFuel, scheme.sigmaBurn,
		scheme.etaBurn, scheme.initAlpha, t0, scheme.nodePrecision,
	)
	nodes[burnerName] = burner
	var freeTurbine = constructive.NewFreeTurbineNode(
		scheme.etaTFree, scheme.lambdaOutFree, scheme.lambdaOutFree, func(constructive.TurbineNode) float64 {
			return scheme.dMassRateRelFree
		},
	)
	nodes[freeTurbineName] = freeTurbine
	var compressorTurbinePowerSink = sink.NewPowerSinkNode()
	nodes[compressorTurbinePowerSinkName] = compressorTurbinePowerSink
	var freeTurbinePowerSink = sink.NewPowerSinkNode()
	nodes[freeTurbinePowerSinkName] = freeTurbinePowerSink
	var pressureLossNode = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["loss"] = pressureLossNode
	var assembler = helper.NewGasStateAssemblerNode()
	nodes["assembler"] =assembler
	var disassembler = helper.NewGasStateDisassemblerNode()
	nodes["disassembler"] = disassembler
	var tSink = sink.NewTemperatureSinkNode()
	nodes["tSink1"] = tSink
	var mSink = sink.NewMassRateRelSinkNode()
	nodes["mSink"] = mSink
	var gSink = sink.NewGasSinkNode()
	nodes["gSink"] = gSink
	var hub = helper.NewHubNode(states.NewPressurePortState(1e5))
	nodes["hub"] = hub
	var regenerator = constructive.NewRegeneratorNode(0.9, 0.05, constructive.SigmaByColdSide)
	nodes["regenerator"] = regenerator
	var outflow = constructive.NewPressureLossNode(pressureLossSigma)
	nodes["outflow"] = outflow
	var breaker = helper.NewCycleBreakerNode(states.StandardAtmosphereState())
	nodes["breaker"] = breaker

	core.Link(compressor.ComplexGasInput(), inputGasSource.ComplexGasOutput())
	core.Link(compressor.ComplexGasOutput(), regenerator.ColdInput())
	core.Link(regenerator.ColdOutput(), burner.ComplexGasInput())
	core.Link(burner.ComplexGasOutput(), compressorTurbine.ComplexGasInput())
	core.Link(compressorTurbine.ComplexGasOutput(), pressureLossNode.ComplexGasInput())
	core.Link(pressureLossNode.ComplexGasOutput(), freeTurbine.ComplexGasInput())

	core.Link(outputGasSource.ComplexGasOutput(), outflow.ComplexGasOutput())
	core.Link(outflow.ComplexGasInput(), disassembler.ComplexGasPort())
	core.Link(disassembler.TemperaturePort(), tSink.TemperatureInput())
	core.Link(disassembler.GasPort(), gSink.GasInput())
	core.Link(disassembler.MassRateRelPort(), mSink.MassRateRelInput())

	core.Link(disassembler.PressurePort(), hub.Inlet())
	core.Link(hub.Outlet1(), freeTurbine.PressureOutput())

	core.Link(hub.Outlet2(), assembler.PressurePort())
	core.Link(freeTurbine.TemperatureOutput(), assembler.TemperaturePort())
	core.Link(freeTurbine.MassRateRelOutput(), assembler.MassRateRelPort())
	core.Link(freeTurbine.GasOutput(), assembler.GasPort())

	core.Link(assembler.ComplexGasPort(), breaker.DataSourcePort())
	core.Link(breaker.UpdatePort(), regenerator.HotInput())
	core.Link(regenerator.HotOutput(), gasSink.ComplexGasInput())

	core.Link(compressor.PowerOutput(), compressorTurbine.PowerInput())
	core.Link(compressorTurbine.PowerOutput(), compressorTurbinePowerSink.PowerInput())
	core.Link(freeTurbine.PowerOutput(), freeTurbinePowerSink.PowerInput())

	var network = core.NewNetwork(nodes)
}
