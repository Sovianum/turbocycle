package free3n

import (
	"github.com/Sovianum/turbocycle/core/graph"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/material/gases"
)

type ThreeShaftCoolFreeScheme interface {
	ThreeShaftFreeScheme
	Cooler() c.CoolerNode
}

func NewThreeShaftCoolFreeScheme(
	gas gases.Gas, tAtm, pAtm, tGas float64,
	mpc c.ParametricCompressorNode, mpcPipe c.PressureLossNode,
	mpt c.ParametricTurbineNode, mptPipe c.PressureLossNode,
	mpEtaM float64,

	hpc c.ParametricCompressorNode, hpcPipe c.PressureLossNode,
	hpt c.ParametricTurbineNode, hptPipe c.PressureLossNode,
	hpEtaM float64,

	lpt c.ParametricTurbineNode, lptPipe c.PressureLossNode,
	burner c.ParametricBurnerNode, payload c.Payload,
	cooler c.CoolerNode,
) ThreeShaftCoolFreeScheme {
	var root = NewThreeShaftFreeScheme(
		gas, tAtm, pAtm, tGas,
		mpc, mpcPipe, mpt, mptPipe, mpEtaM,
		hpc, hpcPipe, hpt, hptPipe, hpEtaM,
		lpt, lptPipe, burner, payload,
	)

	var result = &threeShaftCoolFreeScheme{
		threeShaftFreeScheme: root.(*threeShaftFreeScheme),
		cooler:               cooler,
	}
	result.relink()

	return result
}

type threeShaftCoolFreeScheme struct {
	*threeShaftFreeScheme
	cooler c.CoolerNode
}

func (scheme *threeShaftCoolFreeScheme) Cooler() c.CoolerNode {
	return scheme.cooler
}

func (scheme *threeShaftCoolFreeScheme) relink() {
	graph.LinkAll(
		[]graph.Port{
			scheme.mpcPipe.GasOutput(), scheme.mpcPipe.TemperatureOutput(),
			scheme.mpcPipe.PressureOutput(), scheme.mpcPipe.MassRateOutput(),
		},
		[]graph.Port{
			scheme.cooler.GasInput(), scheme.cooler.TemperatureInput(),
			scheme.cooler.PressureInput(), scheme.cooler.MassRateInput(),
		},
	)
	graph.LinkAll(
		[]graph.Port{
			scheme.cooler.GasOutput(), scheme.cooler.TemperatureOutput(),
			scheme.cooler.PressureOutput(),
		},
		[]graph.Port{
			scheme.hpShaft.Compressor.GasInput(), scheme.hpShaft.Compressor.TemperatureInput(),
			scheme.hpShaft.Compressor.PressureInput(),
		},
	)
	sink.SinkAll(
		scheme.cooler.MassRateOutput(), scheme.hpShaft.Compressor.MassRateInput(),
	)
}
