package nodes

import (
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"math"
)

// TODO add collector port
type compressorNode struct {
	ports     core.PortsType
	EtaAd     float64
	Precision float64
	PiStag    float64
}

func NewCompressorNode(etaAd, piStag, precision float64) *compressorNode {
	var result = &compressorNode{
		ports:     make(core.PortsType),
		EtaAd:     etaAd,
		Precision: precision,
		PiStag:    piStag,
	}

	result.ports[gasInput] = core.NewPort()
	result.ports[gasInput].SetInnerNode(result)

	result.ports[gasOutput] = core.NewPort()
	result.ports[gasOutput].SetInnerNode(result)

	result.ports[powerOutput] = core.NewPort()
	result.ports[powerOutput].SetInnerNode(result)

	return result
}

func (node *compressorNode) GetRequirePortTags() []string {
	return []string{gasInput}
}

func (node *compressorNode) GetUpdatePortTags() []string {
	return []string{gasOutput, powerOutput}
}

func (node *compressorNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *compressorNode) GasInput() *core.Port {
	return node.gasInput()
}

func (node *compressorNode) GasOutput() *core.Port {
	return node.gasOutput()
}

func (node *compressorNode) PowerOutput() *core.Port {
	return node.powerOutput()
}

func (node *compressorNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *compressorNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *compressorNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *compressorNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *compressorNode) LSpecific() float64 {
	return node.lSpecific()
}

func (node *compressorNode) Process() error {
	if node.PiStag <= 1 {
		return errors.New(fmt.Sprintf("Invalid piStag = %f", node.PiStag))
	}

	var pStagOut = node.pStagIn() * node.PiStag
	var tStagOut = node.getTStagOut(node.PiStag, node.tStagIn(), node.tStagIn())

	var gasState = node.GasInput().GetState().(states.GasPortState)
	gasState.TStag = tStagOut
	gasState.PStag = pStagOut

	node.gasOutput().SetState(gasState)

	node.powerOutput().SetState(states.NewPowerPortState(-node.lSpecific()))
	// TODO add and set collector port

	return nil
}

func (node *compressorNode) lSpecific() float64 {
	var cpMean = gases.CpMean(node.gas(), node.tStagIn(), node.tStagOut(), defaultN)
	return cpMean * (node.tStagOut() - node.tStagIn())
}

func (node *compressorNode) getTStagOut(piCStag, tStagIn, tStagOutInit float64) float64 {
	var k = gases.K(node.gas(), tStagIn)
	var x = math.Pow(piCStag, (k-1)/k)

	var tOutCurr = tStagIn * (1 + (x-1)/node.EtaAd)
	var tOutNext = node.tStagOutNewFunc(piCStag, tStagIn, tStagOutInit)

	for !common.Converged(tOutCurr, tOutNext, node.Precision) {
		tOutCurr = tOutNext
		tOutNext = node.tStagOutNewFunc(piCStag, tStagIn, tStagOutInit)
	}

	return tOutNext
}

func (node *compressorNode) tStagOutNewFunc(piCStag, tStagIn, tStagOutCurr float64) float64 {
	var x = node.xFunc(piCStag, tStagIn, tStagOutCurr)
	return tStagIn * (1 + (x-1)/node.EtaAd)
}

func (node *compressorNode) xFunc(piCStag, tStagIn, tStagOut float64) float64 {
	var k = gases.KMean(node.gas(), tStagIn, tStagOut, defaultN)
	return math.Pow(piCStag, (k-1)/k)
}

func (node *compressorNode) tStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).TStag
}

func (node *compressorNode) tStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).TStag
}

func (node *compressorNode) pStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).PStag
}

func (node *compressorNode) pStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).PStag
}

func (node *compressorNode) gas() gases.Gas {
	return node.ports[gasInput].GetState().(states.GasPortState).Gas
}

func (node *compressorNode) gasInput() *core.Port {
	return node.ports[gasInput]
}

func (node *compressorNode) gasOutput() *core.Port {
	return node.ports[gasOutput]
}

func (node *compressorNode) powerOutput() *core.Port {
	return node.ports[powerOutput]
}
