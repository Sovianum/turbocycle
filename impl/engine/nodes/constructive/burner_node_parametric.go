package constructive

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
)

type ParametricBurnerNode interface {
	BurnerNode
	SetFuelRateRel(fuelRateRel float64)
}

func NewParametricBurnerFromProto(
	b BurnerNode, lambdaIn0, massRate0, precision, relaxCoef float64, iterLimit int,
) ParametricBurnerNode {
	pBurn := NewParametricBurnerNode(
		b.Fuel(), b.TFuel(), b.T0(), b.Eta(),
		lambdaIn0, b.PStagIn(), b.TStagIn(),
		massRate0, b.FuelRateRel(), precision, relaxCoef, iterLimit,
		func(lambda float64) float64 {
			return b.Sigma() // todo make something more precise
		},
	)

	graph.CopyAll(
		[]graph.Port{
			b.GasInput(), b.TemperatureInput(), b.PressureInput(), b.MassRateInput(),
			b.GasOutput(), b.TemperatureOutput(), b.PressureOutput(), b.MassRateOutput(),
		},
		[]graph.Port{
			pBurn.GasInput(), pBurn.TemperatureInput(), pBurn.PressureInput(), pBurn.MassRateInput(),
			pBurn.GasOutput(), pBurn.TemperatureOutput(), pBurn.PressureOutput(), pBurn.MassRateOutput(),
		},
	)
	return pBurn
}

func NewParametricBurnerNode(
	fuel fuel.GasFuel, tFuel, t0, etaBurn,
	lambdaIn0, pStagIn0, tStagIn0, massRateIn0, fuelMassRateRel0,
	precision, relaxCoef float64, iterLimit int, sigmaFunc func(lambda float64) float64,
) ParametricBurnerNode {
	var result = &parametricBurnerNode{
		fuelMassRateRel: fuelMassRateRel0,

		lambdaIn0:   lambdaIn0,
		pStagIn0:    pStagIn0,
		tStagIn0:    tStagIn0,
		massRateIn0: massRateIn0,

		precision: precision,
		relaxCoef: relaxCoef,
		iterLimit: iterLimit,

		sigmaFunc: sigmaFunc,
	}
	result.baseBurner = newBaseBurner(result, fuel, etaBurn, tFuel, t0, precision)

	return result
}

type parametricBurnerNode struct {
	*baseBurner

	fuelMassRateRel float64

	lambdaIn0   float64
	pStagIn0    float64
	tStagIn0    float64
	massRateIn0 float64

	sigmaFunc func(lambda float64) float64

	precision float64
	relaxCoef float64
	iterLimit int
}

func (node *parametricBurnerNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ParametricBurner")
}

func (node *parametricBurnerNode) Sigma() float64 {
	return node.sigmaFunc(node.lambdaIn())
}

func (node *parametricBurnerNode) Alpha() float64 {
	return node.alpha()
}

func (node *parametricBurnerNode) FuelRateRel() float64 {
	return node.fuelMassRateRel
}

func (node *parametricBurnerNode) SetFuelRateRel(fuelRateRel float64) {
	node.fuelMassRateRel = fuelRateRel
}

func (node *parametricBurnerNode) Process() error {
	var tStagOut, err = node.tGas()
	if err != nil {
		return err
	}

	var massRateOut = node.massRateInput.GetState().(states.MassRatePortState).MassRate * (1 + node.fuelMassRateRel)

	var sigma = node.sigmaFunc(node.lambdaIn())
	var pStagOut = node.pStagIn() * sigma

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(node.outletGas()), states.NewTemperaturePortState(tStagOut),
			states.NewPressurePortState(pStagOut), states.NewMassRatePortState(massRateOut),
		},
		[]graph.Port{node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput},
	)

	return nil
}

// it is assumed that inlet lambda is low and static density is approximately equal to stagnation one
// the difference of adiabatic parameters of inlet gases on different regimes is also not taken into account
func (node *parametricBurnerNode) lambdaIn() float64 {
	var tFactor = math.Sqrt(node.tStagIn() / node.tStagIn0)
	var pFactor = node.pStagIn0 / node.pStagIn()
	var massRateFactor = node.massRateInput.GetState().(states.MassRatePortState).MassRate / node.massRateIn0

	return tFactor * pFactor * massRateFactor * node.lambdaIn0
}

func (node *parametricBurnerNode) tGas() (float64, error) {
	alphaFunc := func(alpha float64) float64 {
		if alpha <= 1 {
			return alpha
		}
		return 1
	}

	iterFunc := func(tGas float64) (float64, error) {
		cpGas := gases.CpMean(node.outletGas(), tGas, node.t0, nodes.DefaultN)

		tInput := node.tStagIn()
		cpInput := gases.CpMean(node.inletGas(), tInput, node.t0, nodes.DefaultN)

		alpha := node.alpha()

		enom1 := cpInput * (tInput - node.t0)
		enom2 := node.fuelMassRateRel * node.fuel.QLower() * node.etaBurn * alphaFunc(alpha)
		enom3 := node.fuelMassRateRel * node.fuel.Cp(node.t0) * (node.tFuel - node.t0)

		denom := cpGas * (node.fuelMassRateRel + 1)

		return (enom1+enom2+enom3)/denom + node.t0, nil
	}

	return common.SolveIteratively(iterFunc, node.tStagIn(), node.precision, node.relaxCoef, node.iterLimit)
}

func (node *parametricBurnerNode) outletGas() gases.Gas {
	return node.fuel.GetCombustionGas(node.inletGas(), node.alpha())
}

func (node *parametricBurnerNode) inletGas() gases.Gas {
	return node.gasInput.GetState().(states.GasPortState).Gas
}

func (node *parametricBurnerNode) alpha() float64 {
	var gasMassTheory = node.fuel.GasMassTheory(node.inletGas())
	return 1 / (node.fuelMassRateRel * gasMassTheory)
}
