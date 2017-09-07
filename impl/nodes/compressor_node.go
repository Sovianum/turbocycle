package nodes

import (
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"math"
	"github.com/Sovianum/turbocycle/common"
)

// TODO add collector port
type CompressorNode struct {
	ports     core.PortsType
	EtaAd     float64
	Precision float64
	PiStag    float64
}

func NewCompressorNode(etaAd, piStag, precision float64) *CompressorNode {
	var result = &CompressorNode{
		ports:     make(core.PortsType),
		EtaAd:     etaAd,
		Precision: precision,
		PiStag:    piStag,
	}

	result.ports[gasInput] = core.NewPort()
	result.ports[gasInput].SetDest(result)

	result.ports[gasOutput] = core.NewPort()
	result.ports[gasOutput].SetSrc(result)

	result.ports[powerOutput] = core.NewPort()
	result.ports[powerOutput].SetSrc(result)

	return result
}

func (node *CompressorNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *CompressorNode) GasInput() *core.Port {
	return node.gasInput()
}

func (node *CompressorNode) GasOutput() *core.Port {
	return node.gasOutput()
}

func (node *CompressorNode) PowerOutput() *core.Port {
	return node.powerOutput()
}

func (node *CompressorNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *CompressorNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *CompressorNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *CompressorNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *CompressorNode) LSpecific() float64 {
	return node.lSpecific()
}

func (node *CompressorNode) Process() error {
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

func (node *CompressorNode) lSpecific() float64 {
	var cpMean = gases.CpMean(node.gas(), node.tStagIn(), node.tStagOut(), defaultN)
	return cpMean * (node.tStagOut() - node.tStagIn())
}

func (node *CompressorNode) getTStagOut(piCStag, tStagIn, tStagOutInit float64) float64 {
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

func (node *CompressorNode) tStagOutNewFunc(piCStag, tStagIn, tStagOutCurr float64) float64 {
	var x = node.xFunc(piCStag, tStagIn, tStagOutCurr)
	return tStagIn * (1 + (x-1)/node.EtaAd)
}

func (node *CompressorNode) xFunc(piCStag, tStagIn, tStagOut float64) float64 {
	var k = gases.KMean(node.gas(), tStagIn, tStagOut, defaultN)
	return math.Pow(piCStag, (k-1)/k)
}

func (node *CompressorNode) tStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).TStag
}

func (node *CompressorNode) tStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).TStag
}

func (node *CompressorNode) pStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).PStag
}

func (node *CompressorNode) pStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).PStag
}

func (node *CompressorNode) gas() gases.Gas {
	return node.ports[gasInput].GetState().(states.GasPortState).Gas
}

func (node *CompressorNode) gasInput() *core.Port {
	return node.ports[gasInput]
}

func (node *CompressorNode) gasOutput() *core.Port {
	return node.ports[gasOutput]
}

func (node *CompressorNode) powerOutput() *core.Port {
	return node.ports[powerOutput]
}
