package constructive

import (
	"fmt"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
)

type BurnerNode interface {
	graph.Node

	nodes.PressureIn
	nodes.PressureOut
	nodes.TemperatureIn
	nodes.TemperatureOut

	nodes.PressureChannel
	nodes.TemperatureChannel
	nodes.GasChannel
	nodes.MassRateChannel

	Alpha() float64
	FuelRateRel() float64
	Fuel() fuel.GasFuel
	Eta() float64
	Sigma() float64
	T0() float64
	TFuel() float64
}

// while calculating labour function takes massRateRel into account
func FuelMassRate(node BurnerNode) float64 {
	var massRateRel = node.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return node.FuelRateRel() * massRateRel
}

func NewBurnerNode(
	fuel fuel.GasFuel, tgStag, tFuel, sigma, etaBurn, initAlpha, t0, precision, relaxCoef float64, iterLimit int,
) BurnerNode {
	var result = &burnerNode{
		tgStag:    tgStag,
		sigma:     sigma,
		initAlpha: initAlpha,
		precision: precision,
		relaxCoef: relaxCoef,
		iterLimit: iterLimit,
	}
	result.baseBurner = newBaseBurner(result, fuel, etaBurn, tFuel, t0, precision)

	return result
}

type burnerNode struct {
	*baseBurner

	outletGas gases.Gas
	tgStag    float64
	sigma     float64
	initAlpha float64
	alpha     float64
	precision float64
	relaxCoef float64
	iterLimit int
}

func (node *burnerNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Burner")
}

func (node *burnerNode) Sigma() float64 {
	return node.sigma
}

func (node *burnerNode) Alpha() float64 {
	return node.alpha
}

func (node *burnerNode) FuelRateRel() float64 {
	return node.getFuelMassRateRel(node.alpha)
}

func (node *burnerNode) Process() error {
	var fuelMassRateRel, alpha, err = node.getFuelParameters(node.initAlpha)
	if err != nil {
		return err
	}
	node.alpha = alpha

	gasOut := node.outletGas
	tStagOut := node.tgStag
	pStagOut := node.pStagIn() * node.sigma
	massRateRelOut := node.massRateInput.GetState().(states.MassRatePortState).MassRate * (1 + fuelMassRateRel)

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gasOut), states.NewTemperaturePortState(tStagOut),
			states.NewPressurePortState(pStagOut), states.NewMassRatePortState(massRateRelOut),
		},
		[]graph.Port{node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput},
	)

	return nil
}

func (node *burnerNode) getFuelParameters(initAlpha float64) (float64, float64, error) {
	alpha, err := common.SolveIterativelyWithValidation(node.getNextAlpha, func(alpha float64) error {
		if alpha < 0 {
			return fmt.Errorf("invalid alpha: %f", alpha)
		}
		return nil
	}, initAlpha, node.precision, node.relaxCoef, node.iterLimit)
	if err != nil {
		return 0, 0, err
	}

	var fuelMassRateRel = node.getFuelMassRateRel(alpha)
	return fuelMassRateRel, alpha, nil
}

func (node *burnerNode) getNextAlpha(currAlpha float64) (float64, error) {
	var gasMassTheory = node.fuel.GasMassTheory(node.inletGas())
	return 1 / (node.getFuelMassRateRel(currAlpha) * gasMassTheory), nil
}

func (node *burnerNode) getFuelMassRateRel(currAlpha float64) float64 {
	node.outletGas = node.fuel.GetCombustionGas(node.inletGas(), currAlpha)

	num1 := gases.CpMean(node.outletGas, node.tgStag, node.t0, nodes.DefaultN) * (node.tgStag - node.t0)
	num2 := -gases.CpMean(node.inletGas(), node.tStagIn(), node.t0, nodes.DefaultN) * (node.tStagIn() - node.t0)

	denom1 := node.fuel.QLower() * node.etaBurn
	denom2 := -gases.CpMean(node.outletGas, node.tgStag, node.t0, nodes.DefaultN) * (node.tgStag - node.t0)
	denom3 := fuel.CpMean(node.fuel, node.tFuel, node.t0, nodes.DefaultN) * (node.tFuel - node.t0)

	return (num1 + num2) / (denom1 + denom2 + denom3)
}

func (node *burnerNode) inletGas() gases.Gas {
	return node.gasInput.GetState().(states.GasPortState).Gas
}
