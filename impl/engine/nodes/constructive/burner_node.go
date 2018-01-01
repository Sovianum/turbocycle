package constructive

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/go-errors/errors"
)

type BurnerNode interface {
	graph.Node
	nodes.ComplexGasChannel
	nodes.PressureIn
	nodes.PressureOut
	nodes.TemperatureIn
	nodes.TemperatureOut
	Alpha() float64
	GetFuelRateRel() float64
	Fuel() fuel.GasFuel
	Eta() float64
	Sigma() float64
	T0() float64
	TFuel() float64
}

// while calculating labour function takes massRateRel into account
func FuelMassRate(node BurnerNode) float64 {
	var massRateRel = node.ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return node.GetFuelRateRel() * massRateRel
}

func NewBurnerNode(
	fuel fuel.GasFuel, tgStag, tFuel, sigma, etaBurn, initAlpha, t0, precision float64,
) BurnerNode {
	var result = &burnerNode{
		fuel:      fuel,
		tgStag:    tgStag,
		tFuel:     tFuel,
		sigma:     sigma,
		etaBurn:   etaBurn,
		initAlpha: initAlpha,
		t0:        t0,
		precision: precision,
	}

	result.complexGasInput = graph.NewAttachedPort(result)
	result.complexGasOutput = graph.NewAttachedPort(result)

	return result
}

type burnerNode struct {
	graph.BaseNode

	complexGasInput  graph.Port
	complexGasOutput graph.Port

	fuel      fuel.GasFuel
	outletGas gases.Gas
	tgStag    float64
	tFuel     float64
	sigma     float64
	etaBurn   float64
	initAlpha float64
	t0        float64
	precision float64
	alpha     float64
}

func (node *burnerNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Burner")
}

func (node *burnerNode) GetPorts() []graph.Port {
	return []graph.Port{node.complexGasInput, node.complexGasOutput}
}

func (node *burnerNode) GetRequirePorts() []graph.Port {
	return []graph.Port{node.complexGasInput}
}

func (node *burnerNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{node.complexGasOutput}
}

func (node *burnerNode) Fuel() fuel.GasFuel {
	return node.fuel
}

func (node *burnerNode) Eta() float64 {
	return node.etaBurn
}

func (node *burnerNode) Sigma() float64 {
	return node.sigma
}

func (node *burnerNode) T0() float64 {
	return node.t0
}

func (node *burnerNode) TFuel() float64 {
	return node.tFuel
}

func (node *burnerNode) ComplexGasInput() graph.Port {
	return node.complexGasInput
}

func (node *burnerNode) ComplexGasOutput() graph.Port {
	return node.complexGasOutput
}

func (node *burnerNode) Alpha() float64 {
	return node.alpha
}

func (node *burnerNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *burnerNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *burnerNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *burnerNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *burnerNode) GetFuelRateRel() float64 {
	return node.getFuelMassRateRel(node.alpha)
}

func (node *burnerNode) Process() error {
	var fuelMassRateRel, alpha, err = node.getFuelParameters(node.initAlpha)
	if err != nil {
		return err
	}
	node.alpha = alpha

	var gasState = node.complexGasInput.GetState().(states.ComplexGasPortState)
	gasState.Gas = node.outletGas
	gasState.TStag = node.tgStag
	gasState.PStag = node.pStagIn() * node.sigma
	gasState.MassRateRel *= 1 + fuelMassRateRel

	node.complexGasOutput.SetState(gasState)

	return nil
}

func (node *burnerNode) getFuelParameters(initAlpha float64) (float64, float64, error) {
	var currAlpha = initAlpha
	var nextAlpha = node.getNextAlpha(currAlpha)

	for !common.Converged(currAlpha, nextAlpha, node.precision) {
		if math.IsNaN(currAlpha) || math.IsNaN(nextAlpha) {
			return 0, 0, errors.New("failed to converge: try different initial guess")
		}
		currAlpha = nextAlpha
		nextAlpha = node.getNextAlpha(currAlpha)
	}

	var fuelMassRateRel = node.getFuelMassRateRel(nextAlpha)
	return fuelMassRateRel, nextAlpha, nil
}

func (node *burnerNode) getNextAlpha(currAlpha float64) float64 {
	return 1 / (node.getFuelMassRateRel(currAlpha) * node.fuel.AirMassTheory())
}

func (node *burnerNode) getFuelMassRateRel(currAlpha float64) float64 {
	node.outletGas = node.fuel.GetCombustionGas(currAlpha)

	var num1 = gases.CpMean(node.outletGas, node.tgStag, node.t0, nodes.DefaultN) * (node.tgStag - node.t0)
	var num2 = -gases.CpMean(node.inletGas(), node.tStagIn(), node.t0, nodes.DefaultN) * (node.tStagIn() - node.t0)

	var denom1 = node.fuel.QLower() * node.etaBurn
	var denom2 = -gases.CpMean(node.outletGas, node.tgStag, node.t0, nodes.DefaultN) * (node.tgStag - node.t0)
	var denom3 = fuel.CpMean(node.fuel, node.tFuel, node.t0, nodes.DefaultN) * (node.tFuel - node.t0)

	return (num1 + num2) / (denom1 + denom2 + denom3)
}

func (node *burnerNode) inletGas() gases.Gas {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).Gas
}

func (node *burnerNode) tStagIn() float64 {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).TStag
}

func (node *burnerNode) tStagOut() float64 {
	return node.complexGasOutput.GetState().(states.ComplexGasPortState).TStag
}

func (node *burnerNode) pStagIn() float64 {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).PStag
}

func (node *burnerNode) pStagOut() float64 {
	return node.complexGasOutput.GetState().(states.ComplexGasPortState).PStag
}
