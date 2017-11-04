package constructive

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"math"
	"github.com/Sovianum/turbocycle/helpers/gdf"
)

type FreeTurbineNode interface {
	TurbineNode
	nodes.PressureSource
	nodes.TemperatureSource
	nodes.MassRateRelSource
	nodes.GasSource
}

type freeTurbineNode struct {
	ports           core.PortsType
	etaT            float64
	precision       float64
	lambdaOut       float64
	leakMassRateFunc  func(TurbineNode) float64
	coolMasRateRel    func(TurbineNode) float64
	inflowMassRateRel func(TurbineNode) float64
}

func NewFreeTurbineNode(
	etaT, lambdaOut, precision float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(TurbineNode) float64,
) FreeTurbineNode {
	var result = &freeTurbineNode{
		ports:           make(core.PortsType),
		etaT:            etaT,
		precision:       precision,
		lambdaOut:       lambdaOut,
		leakMassRateFunc:  leakMassRateFunc,
		coolMasRateRel:    coolMasRateRel,
		inflowMassRateRel: inflowMassRateRel,
	}

	result.ports[nodes.ComplexGasInput] = core.NewPort()
	result.ports[nodes.ComplexGasInput].SetInnerNode(result)
	result.ports[nodes.ComplexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.PowerOutput] = core.NewPort()
	result.ports[nodes.PowerOutput].SetInnerNode(result)
	result.ports[nodes.PowerOutput].SetState(states.StandardPowerState())

	result.ports[nodes.PressureOutput] = core.NewPort()
	result.ports[nodes.PressureOutput].SetInnerNode(result)
	result.ports[nodes.PressureOutput].SetState(states.NewPressurePortState(1e5)) // todo remove hardcode

	result.ports[nodes.TemperatureOutput] = core.NewPort()
	result.ports[nodes.TemperatureOutput].SetInnerNode(result)
	result.ports[nodes.TemperatureOutput].SetState(states.NewTemperaturePortState(300)) // todo remove hardcode

	result.ports[nodes.GasOutput] = core.NewPort()
	result.ports[nodes.GasOutput].SetInnerNode(result)
	result.ports[nodes.GasOutput].SetState(states.NewGasPortState(gases.GetAir())) // todo remove hardcode

	result.ports[nodes.MassRateRelOutput] = core.NewPort()
	result.ports[nodes.MassRateRelOutput].SetInnerNode(result)
	result.ports[nodes.MassRateRelOutput].SetState(states.NewMassRateRelPortState(1)) // todo remove hardcode

	return result
}

func (node *freeTurbineNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState    core.PortState `json:"gas_input_state"`
		GasOutputState   core.PortState `json:"gas_output_state"`
		PowerOutputState core.PortState `json:"power_output_state"`
		PiStag           float64        `json:"pi_stag"`
		LSpecific        float64        `json:"l_specific"`
		Eta              float64        `json:"eta"`
	}{
		GasInputState: node.complexGasInput().GetState(),
		GasOutputState: states.NewComplexGasPortState(
			node.inputGas(), node.tStagOut(), node.pStagOut(), node.massRateRelOut(),
		),
		PowerOutputState: node.powerOutput().GetState(),
		PiStag:           node.PiTStag(),
		LSpecific:        node.lSpecific(),
		Eta:              node.etaT,
	})
}

func (node *freeTurbineNode) ContextDefined() bool {
	return true
}

func (node *freeTurbineNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasInput:
		return node.complexGasInput(), nil
	case nodes.GasOutput:
		return node.gasOutput(), nil
	case nodes.PowerOutput:
		return node.powerOutput(), nil
	case nodes.PressureOutput:
		return node.pressureOutput(), nil
	case nodes.TemperatureOutput:
		return node.temperatureOutput(), nil
	case nodes.MassRateRelOutput:
		return node.massRateRelOutput(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found in free turbine", tag)
	}
}

func (node *freeTurbineNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.ComplexGasInput, nodes.PressureOutput}, nil
}

func (node *freeTurbineNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.TemperatureOutput, nodes.MassRateRelOutput, nodes.GasOutput, nodes.PowerOutput}, nil
}

func (node *freeTurbineNode) GetPortTags() []string {
	return []string{nodes.ComplexGasInput, nodes.TemperatureOutput, nodes.MassRateRelOutput, nodes.GasOutput, nodes.PowerOutput}
}

func (node *freeTurbineNode) ComplexGasInput() core.Port {
	return node.complexGasInput()
}

func (node *freeTurbineNode) PowerOutput() core.Port {
	return node.powerOutput()
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
	return node.massRateRel()
}

func (node *freeTurbineNode) LeakMassRateRel() float64 {
	return node.leakMassRateFunc(node)
}

func (node *freeTurbineNode) CoolMassRateRel() float64 {
	return node.coolMasRateRel(node)
}

func (node *freeTurbineNode) GetPorts() core.PortsType {
	return node.ports
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
	var gasState = node.complexGasInput().GetState().(states.ComplexGasPortState)

	var tStagOut, err = node.getTStagOut()
	if err != nil {
		return err
	}

	node.temperatureOutput().SetState(states.NewTemperaturePortState(tStagOut))
	node.pressureOutput().SetState(states.NewPressurePortState(node.pStagOut()))
	node.gasOutput().SetState(states.NewGasPortState(gasState.Gas))
	node.massRateRelOutput().SetState(
		states.NewMassRateRelPortState(gasState.MassRateRel * (node.massRateRel())),
	)

	node.powerOutput().SetState(
		states.NewPowerPortState(node.lSpecific()),
	)

	return nil
}

func (node *freeTurbineNode) PressureOutput() core.Port {
	return node.pressureOutput()
}

func (node *freeTurbineNode) TemperatureOutput() core.Port {
	return node.temperatureOutput()
}

func (node *freeTurbineNode) MassRateRelOutput() core.Port {
	return node.massRateRelOutput()
}

func (node *freeTurbineNode) MassRateRelOut() float64 {
	return node.massRateRelOut()
}

func (node *freeTurbineNode) GasOutput() core.Port {
	return node.gasOutput()
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

func (node *freeTurbineNode) massRateRel() float64 {
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
	return node.complexGasInput().GetState().(states.ComplexGasPortState).Gas
}

func (node *freeTurbineNode) tStagIn() float64 {
	return node.complexGasInput().GetState().(states.ComplexGasPortState).TStag
}

func (node *freeTurbineNode) pStagIn() float64 {
	return node.complexGasInput().GetState().(states.ComplexGasPortState).PStag
}

func (node *freeTurbineNode) tStagOut() float64 {
	return node.temperatureOutput().GetState().(states.TemperaturePortState).TStag
}

func (node *freeTurbineNode) pStagOut() float64 {
	return node.pressureOutput().GetState().(states.PressurePortState).PStag
}

func (node *freeTurbineNode) massRateRelOut() float64 {
	return node.massRateRelOutput().GetState().(states.MassRateRelPortState).MassRateRel
}

func (node *freeTurbineNode) temperatureOutput() core.Port {
	return node.ports[nodes.TemperatureOutput]
}

func (node *freeTurbineNode) pressureOutput() core.Port {
	return node.ports[nodes.PressureOutput]
}

func (node *freeTurbineNode) massRateRelOutput() core.Port {
	return node.ports[nodes.MassRateRelOutput]
}

func (node *freeTurbineNode) gasOutput() core.Port {
	return node.ports[nodes.GasOutput]
}

func (node *freeTurbineNode) complexGasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *freeTurbineNode) powerOutput() core.Port {
	return node.ports[nodes.PowerOutput]
}
