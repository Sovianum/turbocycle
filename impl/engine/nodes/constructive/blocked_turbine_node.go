package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type BlockedTurbineNode interface {
	TurbineNode
	nodes.PowerSink
}

func NewBlockedTurbineNode(
	etaT, lambdaOut, precision float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(TurbineNode) float64,
) BlockedTurbineNode {
	var result = &blockedTurbineNode{
		etaT:              etaT,
		precision:         precision,
		lambdaOut:         lambdaOut,
		leakMassRateFunc:  leakMassRateFunc,
		coolMasRateRel:    coolMasRateRel,
		inflowMassRateRel: inflowMassRateRel,
	}

	result.baseBlockedTurbine = NewBaseBlockedTurbine(result, precision)
	return result
}

type blockedTurbineNode struct {
	graph.BaseNode
	*baseBlockedTurbine

	etaT              float64
	precision         float64
	lambdaOut         float64
	leakMassRateFunc  func(TurbineNode) float64
	coolMasRateRel    func(TurbineNode) float64
	inflowMassRateRel func(TurbineNode) float64
}

func (node *blockedTurbineNode) LambdaOut() float64 {
	return node.lambdaOut
}

func (node *blockedTurbineNode) PStatOut() float64 {
	return node.pStatOut(node.lambdaOut)
}

func (node *blockedTurbineNode) TStatOut() float64 {
	return node.tStatOut(node.lambdaOut)
}

func (node *blockedTurbineNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "BlockedTurbine")
}

func (node *blockedTurbineNode) Process() error {
	var tStagOut, err = node.getTStagOut(node.turbineLabour())
	if err != nil {
		return err
	}

	var piTStag = node.piTStag(tStagOut, node.etaT)
	var pi = gdf.Pi(node.lambdaOut, gases.KMean(node.inputGas(), node.tStagIn(), tStagOut, nodes.DefaultN))
	var pStagOut = node.pStagIn() / (piTStag * pi)
	var massRateRelOut = node.massRateInput.GetState().(states.MassRatePortState).MassRate * node.massRateRelFactor()

	node.temperatureOutput.SetState(states.NewTemperaturePortState(tStagOut))
	node.pressureOutput.SetState(states.NewPressurePortState(pStagOut))
	node.gasOutput.SetState(states.NewGasPortState(node.inputGas()))
	node.massRateOutput.SetState(states.NewMassRatePortState(massRateRelOut))

	node.powerOutput.SetState(states.NewPowerPortState(node.turbineLabour())) // TODO maybe need to pass sum of labours

	return nil
}

func (node *blockedTurbineNode) Eta() float64 {
	return node.etaT
}

func (node *blockedTurbineNode) LeakMassRateRel() float64 {
	return node.leakMassRateFunc(node)
}

func (node *blockedTurbineNode) CoolMassRateRel() float64 {
	return node.coolMasRateRel(node)
}

func (node *blockedTurbineNode) PiTStag() float64 {
	return node.piTStag(node.tStagOut(), node.etaT)
}

func (node *blockedTurbineNode) massRateRelFactor() float64 {
	return 1 + node.leakMassRateFunc(node) + node.coolMasRateRel(node) + node.inflowMassRateRel(node)
}
