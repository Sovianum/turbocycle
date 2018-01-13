package constructive

import (
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type BlockedTurbineNode interface {
	StaticTurbineNode
	nodes.PowerSink
	nodes.MassRateSink
}

func NewSimpleBlockedTurbineNode(
	etaT, lambdaOut,
	leakMassRateCoef, coolMassRateCoef, inflowMassRateCoef,
	precision float64,
) BlockedTurbineNode {
	return NewBlockedTurbineNode(
		etaT, lambdaOut, precision,
		func(TurbineNode) float64 {
			return leakMassRateCoef
		},
		func(TurbineNode) float64 {
			return coolMassRateCoef
		},
		func(TurbineNode) float64 {
			return inflowMassRateCoef
		},
	)
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
	result.powerInput = graph.NewAttachedPortWithTag(result, nodes.PowerInputTag)
	result.massRateInput = graph.NewAttachedPortWithTag(result, nodes.MassRateInputTag)
	return result
}

type blockedTurbineNode struct {
	graph.BaseNode
	*baseBlockedTurbine

	powerInput    graph.Port
	massRateInput graph.Port

	etaT              float64
	precision         float64
	lambdaOut         float64
	leakMassRateFunc  func(TurbineNode) float64
	coolMasRateRel    func(TurbineNode) float64
	inflowMassRateRel func(TurbineNode) float64
}

func (node *blockedTurbineNode) LSpecific() float64 {
	return node.turbineLabour()
}

func (node *blockedTurbineNode) LambdaOut() float64 {
	return node.lambdaOut
}

func (node *blockedTurbineNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "BlockedTurbine")
}

func (node *blockedTurbineNode) GetPorts() []graph.Port {
	return append(node.baseBlockedTurbine.GetPorts(), node.massRateInput, node.powerInput)
}

func (node *blockedTurbineNode) GetRequirePorts() ([]graph.Port, error) {
	var ports, err = node.baseBlockedTurbine.GetRequirePorts()
	if err != nil {
		return nil, err
	}
	return append(ports, node.massRateInput, node.powerInput), nil
}

func (node *blockedTurbineNode) Process() error {
	//var tStagOut, err = node.getTStagOut(node.turbineLabour())
	var tStagOut, err = node.getTStagOut()
	if err != nil {
		return err
	}

	var piTStag = node.piTStag(tStagOut, node.etaT)
	var pStagOut = node.pStagIn() / piTStag
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

func (node *blockedTurbineNode) PowerInput() graph.Port {
	return node.powerInput
}

func (node *blockedTurbineNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *blockedTurbineNode) massRateRelFactor() float64 {
	return 1 + node.leakMassRateFunc(node) + node.coolMasRateRel(node) + node.inflowMassRateRel(node)
}

// here it is assumed that pressure drop is calculated by stagnation parameters
func (node *blockedTurbineNode) piTStag(tStagOut, etaT float64) float64 {
	var k = gases.KMean(node.inputGas(), node.tStagIn(), tStagOut, nodes.DefaultN)
	var cp = gases.CpMean(node.inputGas(), node.tStagIn(), tStagOut, nodes.DefaultN)

	var labour = node.turbineLabour()
	return math.Pow(
		1-labour/(cp*node.tStagIn()*etaT),
		k/(1-k),
	)
}

func (node *blockedTurbineNode) getTStagOut() (float64, error) {
	var t0, err = node.getNewTtStag(0.8 * node.tStagIn()) // TODO move 0.8 out of code
	if err != nil {
		return 0, err
	}
	return common.SolveIterativly(node.getNewTtStag, t0, node.precision, nodes.DefaultN)
}

func (node *blockedTurbineNode) getNewTtStag(currTtStag float64) (float64, error) {
	var cp = gases.CpMean(node.inputGas(), node.tStagIn(), currTtStag, nodes.DefaultN)
	var tTStag = node.tStagIn() - node.turbineLabour()/cp
	if math.IsNaN(tTStag) {
		return 0, fmt.Errorf("nan obtained while calculating TtStag")
	}
	return tTStag, nil
}

func (node *blockedTurbineNode) turbineLabour() float64 {
	return -node.powerInput.GetState().(states.PowerPortState).LSpecific
}
