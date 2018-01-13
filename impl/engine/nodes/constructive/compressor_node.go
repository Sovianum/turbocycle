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

func NewCompressorNode(etaPol, piStag, precision float64) CompressorNode {
	var result = &compressorNode{
		etaPol:    etaPol,
		precision: precision,
		piStag:    piStag,
	}
	result.baseCompressor = newBaseCompressor(result, precision)
	result.massRateInput = graph.NewAttachedPortWithTag(result, nodes.MassRateInputTag)
	return result
}

// TODO add collector port
type compressorNode struct {
	*baseCompressor

	massRateInput graph.Port

	etaPol    float64 // politropic efficiency
	precision float64
	piStag    float64
}

// while calculating labour function takes massRateRel into account
func CompressorLabour(node CompressorNode) float64 {
	var massRateRel = node.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return node.LSpecific() * massRateRel
}

func (node *compressorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Compressor")
}

func (node *compressorNode) GetPorts() []graph.Port {
	return append(node.baseCompressor.GetPorts(), node.massRateInput)
}

func (node *compressorNode) GetRequirePorts() ([]graph.Port, error) {
	var ports, err = node.baseCompressor.GetRequirePorts()
	if err != nil {
		return nil, err
	}
	return append(ports, node.massRateInput), nil
}

func (node *compressorNode) PiStag() float64 {
	return node.piStag
}

func (node *compressorNode) Eta() float64 {
	return node.etaAd(node.tStagOut())
}

func (node *compressorNode) EtaPol() float64 {
	return node.etaPol
}

func (node *compressorNode) SetPiStag(piStag float64) {
	node.piStag = piStag
}

func (node *compressorNode) Process() error {
	if node.piStag <= 1 {
		return fmt.Errorf("invalid piStag = %f", node.piStag)
	}

	var tStagOutInit = node.tStagIn()

	var etaAdCurr = node.etaAd(node.tStagIn())
	var tStagOut, err = node.getTStagOut(tStagOutInit, node.piStag, etaAdCurr)
	if err != nil {
		return err
	}
	var etaAdNew = node.etaAd(tStagOut)

	for !common.Converged(etaAdCurr, etaAdNew, node.precision) {
		etaAdCurr = etaAdNew
		tStagOut, err = node.getTStagOut(tStagOutInit, node.piStag, etaAdCurr)
		if err != nil {
			return err
		}
		etaAdNew = node.etaAd(tStagOut)
	}

	var pStagOut = node.pStagIn() * node.piStag

	graph.SetAll(
		[]graph.PortState{
			node.gasInput.GetState(),
			states.NewTemperaturePortState(tStagOut),
			states.NewPressurePortState(pStagOut),
			node.massRateInput.GetState(),
		},
		[]graph.Port{node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput},
	)

	node.powerOutput.SetState(states.NewPowerPortState(-node.lSpecific()))

	return nil
}

func (node *compressorNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *compressorNode) etaAd(tStagOut float64) float64 {
	var k = gases.KMean(node.gas(), node.tStagIn(), tStagOut, nodes.DefaultN)

	var enom = math.Pow(node.piStag, (k-1)/k) - 1
	var denom = math.Pow(node.piStag, (k-1)/(k*node.etaPol)) - 1

	return enom / denom
}
