package constructive

import (
	"errors"
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type CompressorNode interface {
	graph.Node
	nodes.ComplexGasChannel
	nodes.PowerSource
	nodes.PressureIn
	nodes.PressureOut
	nodes.TemperatureIn
	nodes.TemperatureOut
	LSpecific() float64
	PiStag() float64
	Eta() float64
	EtaPol() float64
	SetPiStag(piStag float64)
}

func NewCompressorNode(etaAd, piStag, precision float64) CompressorNode {
	var result = &compressorNode{
		etaPol:    etaAd,
		precision: precision,
		piStag:    piStag,
	}

	result.complexGasInput = graph.NewAttachedPort(result)
	result.complexGasOutput = graph.NewAttachedPort(result)
	result.powerOutput = graph.NewAttachedPort(result)

	return result
}

// TODO add collector port
type compressorNode struct {
	graph.BaseNode

	complexGasInput  graph.Port
	complexGasOutput graph.Port
	powerOutput      graph.Port

	etaPol    float64 // politropic efficiency
	precision float64
	piStag    float64
}

func (node *compressorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Compressor")
}

func (node *compressorNode) GetPorts() []graph.Port {
	return []graph.Port{node.complexGasInput, node.complexGasOutput, node.powerOutput}
}

func (node *compressorNode) GetRequirePorts() []graph.Port {
	return []graph.Port{node.complexGasInput}
}

func (node *compressorNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{node.complexGasOutput, node.powerOutput}
}

// while calculating labour function takes massRateRel into account
func CompressorLabour(node CompressorNode) float64 {
	var massRateRel = node.ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return node.LSpecific() * massRateRel
}

func (node *compressorNode) Process() error {
	if node.piStag <= 1 {
		return fmt.Errorf("Invalid piStag = %f", node.piStag)
	}

	var pStagOut = node.pStagIn() * node.piStag
	var tStagOut, err = node.getTStagOut(node.piStag, node.tStagIn(), node.tStagIn())
	if err != nil {
		return err
	}

	var gasState = node.ComplexGasInput().GetState().(states.ComplexGasPortState)
	gasState.TStag = tStagOut
	gasState.PStag = pStagOut

	node.complexGasOutput.SetState(gasState)

	node.powerOutput.SetState(states.NewPowerPortState(-node.lSpecific()))
	// TODO add and set collector port

	return nil
}

func (node *compressorNode) ComplexGasInput() graph.Port {
	return node.complexGasInput
}

func (node *compressorNode) ComplexGasOutput() graph.Port {
	return node.complexGasOutput
}

func (node *compressorNode) PowerOutput() graph.Port {
	return node.powerOutput
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

func (node *compressorNode) PiStag() float64 {
	return node.piStag
}

func (node *compressorNode) Eta() float64 {
	return node.etaAd(node.piStag, node.tStagIn(), node.tStagOut())
}

func (node *compressorNode) EtaPol() float64 {
	return node.etaPol
}

func (node *compressorNode) SetPiStag(piStag float64) {
	node.piStag = piStag
}

func (node *compressorNode) lSpecific() float64 {
	var cpMean = gases.CpMean(node.gas(), node.tStagIn(), node.tStagOut(), nodes.DefaultN)
	return cpMean * (node.tStagOut() - node.tStagIn())
}

func (node *compressorNode) getTStagOut(piCStag, tStagIn, tStagOutInit float64) (float64, error) {
	var k = gases.K(node.gas(), tStagIn)
	var x = math.Pow(piCStag, (k-1)/k)

	var tOutCurr = tStagIn * (1 + (x-1)/node.etaPol)
	var tOutNext = node.tStagOutNewFunc(piCStag, tStagIn, tStagOutInit)

	for !common.Converged(tOutCurr, tOutNext, node.precision) {
		if math.IsNaN(tOutCurr) || math.IsNaN(tOutNext) {
			return 0, errors.New("failed to converge: try another initial guess")
		}
		tOutCurr = tOutNext
		tOutNext = node.tStagOutNewFunc(piCStag, tStagIn, tStagOutInit)
	}

	return tOutNext, nil
}

func (node *compressorNode) tStagOutNewFunc(piCStag, tStagIn, tStagOutCurr float64) float64 {
	var x = node.xFunc(piCStag, tStagIn, tStagOutCurr)
	var etaAd = node.etaAd(piCStag, tStagIn, tStagOutCurr)

	return tStagIn * (1 + (x-1)/etaAd)
}

func (node *compressorNode) etaAd(piCStag, tStagIn, tStagOut float64) float64 {
	var k = gases.KMean(node.gas(), tStagIn, tStagOut, nodes.DefaultN)

	var enom = math.Pow(piCStag, (k-1)/k) - 1
	var denom = math.Pow(piCStag, (k-1)/(k*node.etaPol)) - 1

	return enom / denom
}

func (node *compressorNode) xFunc(piCStag, tStagIn, tStagOut float64) float64 {
	var k = gases.KMean(node.gas(), tStagIn, tStagOut, nodes.DefaultN)
	return math.Pow(piCStag, (k-1)/k)
}

func (node *compressorNode) tStagIn() float64 {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).TStag
}

func (node *compressorNode) tStagOut() float64 {
	return node.complexGasOutput.GetState().(states.ComplexGasPortState).TStag
}

func (node *compressorNode) pStagIn() float64 {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).PStag
}

func (node *compressorNode) pStagOut() float64 {
	return node.complexGasOutput.GetState().(states.ComplexGasPortState).PStag
}

func (node *compressorNode) gas() gases.Gas {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).Gas
}
