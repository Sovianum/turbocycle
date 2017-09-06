package nodes

import (
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"math"
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
	return node.ports[gasInput]
}

func (node *CompressorNode) GasOutput() *core.Port {
	return node.ports[gasOutput]
}

func (node *CompressorNode) PowerOutput() *core.Port {
	return node.ports[powerOutput]
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

	node.GasOutput().SetState(gasState)

	node.PowerOutput().SetState(states.NewPowerPortState(-node.lSpecific()))
	// TODO add and set collector port

	return nil
}

func (node *CompressorNode) getTStagOut(piCStag, tStagIn, tStagOutInit float64) float64 {
	var xFunc = func(piCStag, tStagIn, tStagOut float64) float64 {
		var k = gases.KMean(node.gas(), tStagIn, tStagOut, defaultN)
		return math.Pow(piCStag, (k-1)/k)
	}

	var tStagOutNewFunc = func(piCStag, tStagIn, tStagOutCurr float64) float64 {
		var x = xFunc(piCStag, tStagIn, tStagOutCurr)
		return tStagIn * (1 + (x-1)/node.EtaAd)
	}

	var tOutInitFunc = func(piCStag, tStagIn float64) float64 {
		var k = gases.K(node.gas(), tStagIn)
		var x = math.Pow(piCStag, (k-1)/k)
		return tStagIn * (1 + (x-1)/node.EtaAd)
	}

	var isValid = func(tCurr, tNext float64) bool {
		return math.Abs(tCurr-tNext)/tCurr <= node.Precision
	}

	var tOutCurr = tOutInitFunc(piCStag, tStagIn)
	var tOutNext = tStagOutNewFunc(piCStag, tStagIn, tStagOutInit)

	for !isValid(tOutCurr, tOutNext) {
		tOutCurr = tOutNext
		tOutNext = tStagOutNewFunc(piCStag, tStagIn, tStagOutInit)
	}

	return tOutNext
}

func (node *CompressorNode) tStagIn() float64 {
	return node.GasInput().GetState().(states.GasPortState).TStag
}

func (node *CompressorNode) tStagOut() float64 {
	return node.GasOutput().GetState().(states.GasPortState).TStag
}

func (node *CompressorNode) pStagIn() float64 {
	return node.GasInput().GetState().(states.GasPortState).PStag
}

func (node *CompressorNode) pStagOut() float64 {
	return node.GasOutput().GetState().(states.GasPortState).PStag
}

func (node *CompressorNode) gas() gases.Gas {
	return node.ports[gasInput].GetState().(states.GasPortState).Gas
}

func (node *CompressorNode) lSpecific() float64 {
	var cpMean = gases.CpMean(node.gas(), node.tStagIn(), node.tStagOut(), defaultN)
	return cpMean * (node.tStagOut() - node.tStagIn())
}
