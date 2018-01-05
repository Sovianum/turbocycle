package constructive

import (
	"errors"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type FreeTurbineNode interface {
	StaticTurbineNode
	nodes.MassRateSink
}

func NewFreeTurbineNode(
	etaT, lambdaOut, precision float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(TurbineNode) float64,
) FreeTurbineNode {
	var result = &freeTurbineNode{
		etaT:              etaT,
		precision:         precision,
		lambdaOut:         lambdaOut,
		leakMassRateFunc:  leakMassRateFunc,
		coolMasRateRel:    coolMasRateRel,
		inflowMassRateRel: inflowMassRateRel,
	}

	graph.AttachAllPorts(
		result,
		&result.powerOutput,
		&result.temperatureInput, &result.pressureInput, &result.gasInput, &result.massRateInput,
		&result.temperatureOutput, &result.pressureOutput, &result.gasOutput, &result.massRateOutput,
	)

	return result
}

type freeTurbineNode struct {
	graph.BaseNode

	powerOutput graph.Port

	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port
	massRateInput    graph.Port

	pressureOutput    graph.Port
	temperatureOutput graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port

	etaT              float64
	precision         float64
	lambdaOut         float64
	leakMassRateFunc  func(TurbineNode) float64
	coolMasRateRel    func(TurbineNode) float64
	inflowMassRateRel func(TurbineNode) float64
}

func (node *freeTurbineNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "FreeTurbine")
}

func (node *freeTurbineNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.powerOutput,
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *freeTurbineNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput, node.pressureOutput,
	}
}

func (node *freeTurbineNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.powerOutput, node.temperatureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *freeTurbineNode) PowerOutput() graph.Port {
	return node.powerOutput
}

func (node *freeTurbineNode) InputGas() gases.Gas {
	return node.inputGas()
}

func (node *freeTurbineNode) LambdaOut() float64 {
	return node.lambdaOut
}

func (node *freeTurbineNode) Eta() float64 {
	return node.etaT
}

func (node *freeTurbineNode) PStatOut() float64 {
	return node.pStatOut()
}

func (node *freeTurbineNode) TStatOut() float64 {
	return node.tStatOut()
}

func (node *freeTurbineNode) MassRateRel() float64 {
	return node.massRateOutput.GetState().(states.MassRatePortState).MassRate
}

func (node *freeTurbineNode) LeakMassRateRel() float64 {
	return node.leakMassRateFunc(node)
}

func (node *freeTurbineNode) CoolMassRateRel() float64 {
	return node.coolMasRateRel(node)
}

func (node *freeTurbineNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *freeTurbineNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *freeTurbineNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *freeTurbineNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *freeTurbineNode) PiTStag() float64 {
	return node.piTStag()
}

func (node *freeTurbineNode) LSpecific() float64 {
	return node.lSpecific()
}

func (node *freeTurbineNode) Process() error {
	//var gasState = node.complexGasInput.GetState().(states.ComplexGasPortState)

	var tStagOut, err = node.getTStagOut()
	if err != nil {
		return err
	}

	node.temperatureOutput.SetState(states.NewTemperaturePortState(tStagOut))
	node.pressureOutput.SetState(states.NewPressurePortState(node.pStagOut()))
	node.gasOutput.SetState(states.NewGasPortState(node.inputGas()))
	node.massRateOutput.SetState(
		states.NewMassRatePortState(node.massRateInput.GetState().(states.MassRatePortState).MassRate *
			node.massRateRelFactor()),
	)

	node.powerOutput.SetState(
		states.NewPowerPortState(node.lSpecific()),
	)

	return nil
}

func (node *freeTurbineNode) GasInput() graph.Port {
	return node.gasInput
}

func (node *freeTurbineNode) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *freeTurbineNode) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *freeTurbineNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *freeTurbineNode) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *freeTurbineNode) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *freeTurbineNode) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *freeTurbineNode) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *freeTurbineNode) lSpecific() float64 {
	return gases.CpMean(node.inputGas(), node.tStagIn(), node.tStagOut(), nodes.DefaultN) * (node.tStagIn() - node.tStagOut())
}

func (node *freeTurbineNode) tStatOut() float64 {
	var tStagOut = node.tStagOut()
	var k = gases.K(node.inputGas(), tStagOut)
	return tStagOut * gdf.Tau(node.lambdaOut, k)
}

func (node *freeTurbineNode) pStatOut() float64 {
	var pStagOut = node.pStagOut()
	var tStagOut = node.tStagOut()
	var k = gases.K(node.inputGas(), tStagOut)
	return pStagOut * gdf.Tau(node.lambdaOut, k)
}

func (node *freeTurbineNode) massRateRelFactor() float64 {
	return 1 + node.leakMassRateFunc(node) + node.coolMasRateRel(node) + node.inflowMassRateRel(node)
}

func (node *freeTurbineNode) getTStagOut() (float64, error) {
	var tStagOutCurr = node.tStagOutNext(
		node.pStagIn(), node.pStagOut(), node.tStagIn(), node.tStagIn(),
	)
	var tStagOutNext = node.tStagOutNext(
		node.pStagIn(), node.pStagOut(), node.tStagIn(), tStagOutCurr,
	)

	for !common.Converged(tStagOutCurr, tStagOutNext, node.precision) {
		if math.IsNaN(tStagOutCurr) || math.IsNaN(tStagOutNext) {
			return 0, errors.New("failed to converge: try different initial guess")
		}
		tStagOutCurr = tStagOutNext
		node.tStagOutNext(
			node.pStagIn(), node.pStagOut(), node.tStagIn(), tStagOutCurr,
		)
	}

	return tStagOutNext, nil
}

func (node *freeTurbineNode) tStagOutNext(pStagIn, pStagOut, tStagIn, tStagOutCurr float64) float64 {
	var k = gases.KMean(node.inputGas(), tStagIn, tStagOutCurr, nodes.DefaultN)
	var piTStag = pStagIn / pStagOut
	var piT = piTStag / gdf.Pi(node.lambdaOut, gases.K(node.InputGas(), tStagOutCurr))
	var x = math.Pow(piT, (1-k)/k)

	return tStagIn * (1 - (1-x)*node.etaT)
}

func (node *freeTurbineNode) piTStag() float64 {
	return node.pStagIn() / node.pStagOut()
}

func (node *freeTurbineNode) inputGas() gases.Gas {
	return node.gasInput.GetState().(states.GasPortState).Gas
}

func (node *freeTurbineNode) tStagIn() float64 {
	return node.temperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *freeTurbineNode) pStagIn() float64 {
	return node.pressureInput.GetState().(states.PressurePortState).PStag
}

func (node *freeTurbineNode) tStagOut() float64 {
	return node.temperatureOutput.GetState().(states.TemperaturePortState).TStag
}

func (node *freeTurbineNode) pStagOut() float64 {
	return node.pressureOutput.GetState().(states.PressurePortState).PStag
}
