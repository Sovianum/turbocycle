package free3n

import (
	"github.com/Sovianum/turbocycle/core/graph"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/material/gases"
)

type ThreeShaftBurnCoolFreeScheme interface {
	ThreeShaftBurnFreeScheme
	Cooler() c.CoolerNode
}

func NewThreeShaftBurnCoolFreeScheme(
	gas gases.Gas, tAtm, pAtm, tGas, tGasMid float64,
	mpc c.ParametricCompressorNode, mpcPipe c.PressureLossNode,
	mpt c.ParametricTurbineNode, mptPipe c.PressureLossNode,
	mpEtaM float64,

	hpc c.ParametricCompressorNode, hpcPipe c.PressureLossNode,
	hpt c.ParametricTurbineNode, hptPipe c.PressureLossNode,
	hpEtaM float64,

	lpt c.ParametricTurbineNode, lptPipe c.PressureLossNode,
	burner c.ParametricBurnerNode, payload c.Payload,
	cooler c.CoolerNode, midBurner c.ParametricBurnerNode,
) ThreeShaftBurnCoolFreeScheme {
	var root = NewThreeShaftBurnFreeScheme(
		gas, tAtm, pAtm, tGas, tGasMid,
		mpc, mpcPipe, mpt, mptPipe, mpEtaM,
		hpc, hpcPipe, hpt, hptPipe, hpEtaM,
		lpt, lptPipe, burner, payload, midBurner,
	)

	var result = &threeShaftBurnCoolFreeScheme{
		threeShaftBurnFreeScheme: root.(*threeShaftBurnFreeScheme),
		cooler: cooler,
	}
	result.relink()

	return result
}

type threeShaftBurnCoolFreeScheme struct {
	*threeShaftBurnFreeScheme
	cooler c.CoolerNode
}

func (scheme *threeShaftBurnCoolFreeScheme) Cooler() c.CoolerNode {
	return scheme.cooler
}

func (scheme *threeShaftBurnCoolFreeScheme) GetNetwork() (graph.Network, error) {
	return graph.NewNetwork(append(scheme.getNodes(), scheme.cooler))
}

func (scheme *threeShaftBurnCoolFreeScheme) relink() {
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
