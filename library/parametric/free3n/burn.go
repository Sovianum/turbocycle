package free3n

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/variator"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive/utils"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/material/gases"
)

type ThreeShaftBurnFreeScheme interface {
	ThreeShaftFreeScheme
	MidBurner() c.ParametricBurnerNode
}

func NewThreeShaftBurnFreeScheme(
	gas gases.Gas, tAtm, pAtm, tGas, tGasMid float64,
	mpc c.ParametricCompressorNode, mpcPipe c.PressureLossNode,
	mpt c.ParametricTurbineNode, mptPipe c.PressureLossNode,
	mpEtaM float64,

	hpc c.ParametricCompressorNode, hpcPipe c.PressureLossNode,
	hpt c.ParametricTurbineNode, hptPipe c.PressureLossNode,
	hpEtaM float64,

	lpt c.ParametricTurbineNode, lptPipe c.PressureLossNode,
	burner c.ParametricBurnerNode, payload c.Payload,
	midBurner c.ParametricBurnerNode,
) ThreeShaftBurnFreeScheme {
	var root = NewThreeShaftFreeScheme(
		gas, tAtm, pAtm, tGas,
		mpc, mpcPipe, mpt, mptPipe, mpEtaM,
		hpc, hpcPipe, hpt, hptPipe, hpEtaM,
		lpt, lptPipe, burner, payload,
	)

	var result = &threeShaftBurnFreeScheme{
		threeShaftFreeScheme: root.(*threeShaftFreeScheme),
		midBurner:            midBurner,
		midBurnerTSource:     source.NewTemperatureSourceNode(tGasMid),
	}
	result.relink()
	result.resetEquations()

	return result
}

type threeShaftBurnFreeScheme struct {
	*threeShaftFreeScheme

	midBurner        c.ParametricBurnerNode
	midBurnerTSource source.TemperatureSourceNode

	midBurnerEq graph.ReduceNode
}

func (scheme *threeShaftBurnFreeScheme) MidBurner() c.ParametricBurnerNode {
	return scheme.midBurner
}

func (scheme *threeShaftBurnFreeScheme) GetNetwork() (graph.Network, error) {
	return graph.NewNetwork(scheme.getNodes())
}

func (scheme *threeShaftBurnFreeScheme) getNodes() []graph.Node {
	return append(
		scheme.threeShaftFreeScheme.getNodes(),
		scheme.midBurner, scheme.midBurnerEq, scheme.midBurnerTSource,
	)
}

func (scheme *threeShaftBurnFreeScheme) relink() {
	graph.LinkAll(
		[]graph.Port{
			scheme.hptPipe.GasOutput(), scheme.hptPipe.TemperatureOutput(),
			scheme.hptPipe.PressureOutput(), scheme.hptPipe.MassRateOutput(),
		},
		[]graph.Port{
			scheme.midBurner.GasInput(), scheme.midBurner.TemperatureInput(),
			scheme.midBurner.PressureInput(), scheme.midBurner.MassRateInput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			scheme.midBurner.GasOutput(), scheme.midBurner.TemperatureOutput(),
			scheme.midBurner.PressureOutput(),
		},
		[]graph.Port{
			scheme.mpShaft.Turbine.GasInput(), scheme.mpShaft.Turbine.TemperatureInput(),
			scheme.mpShaft.Turbine.PressureInput(),
		},
	)
	sink.SinkAll(
		scheme.midBurner.MassRateOutput(), scheme.mpShaft.Turbine.MassRateInput(),
		scheme.midBurnerTSource.TemperatureOutput(),
	)
}

func (scheme *threeShaftBurnFreeScheme) resetEquations() {
	scheme.variators = append(
		scheme.variators,
		variator.FromCallables(scheme.midBurner.FuelRateRel, scheme.midBurner.SetFuelRateRel),
	)
	scheme.midBurnerEq = utils.NewEquality(
		graph.NewWeakPort(scheme.midBurner.TemperatureOutput()),
		graph.NewWeakPort(scheme.midBurnerTSource.TemperatureOutput()),
	)
	scheme.assembler.AddInputPorts(scheme.midBurnerEq.OutputPort())
}
