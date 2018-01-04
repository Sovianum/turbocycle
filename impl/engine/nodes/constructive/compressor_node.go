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

	graph.AttachAllPorts(
		result,
		&result.powerOutput,
		&result.gasInput, &result.temperatureInput, &result.pressureInput, &result.massRateInput,
		&result.gasOutput, &result.temperatureOutput, &result.pressureOutput, &result.massRateOutput,
	)

	return result
}

// TODO add collector port
type compressorNode struct {
	graph.BaseNode

	gasInput         graph.Port
	temperatureInput graph.Port
	pressureInput    graph.Port
	massRateInput    graph.Port

	gasOutput         graph.Port
	temperatureOutput graph.Port
	pressureOutput    graph.Port
	massRateOutput    graph.Port

	powerOutput graph.Port

	etaPol    float64 // politropic efficiency
	precision float64
	piStag    float64
}

// while calculating labour function takes massRateRel into account
func CompressorLabour(node CompressorNode) float64 {
	var massRateRel = node.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return node.LSpecific() * massRateRel
}

func (node *compressorNode) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *compressorNode) GasInput() graph.Port {
	return node.gasInput
}

func (node *compressorNode) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *compressorNode) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *compressorNode) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *compressorNode) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *compressorNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *compressorNode) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *compressorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Compressor")
}

func (node *compressorNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.powerOutput,
		node.gasInput, node.temperatureInput, node.pressureInput, node.massRateInput,
		node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
	}
}

func (node *compressorNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.gasInput, node.temperatureInput, node.pressureInput, node.massRateInput,
	}
}

func (node *compressorNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.powerOutput,
		node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
	}
}

func (node *compressorNode) Process() error {
	if node.piStag <= 1 {
		return fmt.Errorf("invalid piStag = %f", node.piStag)
	}

	var pStagOut = node.pStagIn() * node.piStag
	var tStagOut, err = node.getTStagOut(node.piStag, node.tStagIn(), node.tStagIn())
	if err != nil {
		return err
	}

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
	// TODO add and set collector port

	return nil
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
	return node.temperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *compressorNode) tStagOut() float64 {
	return node.temperatureOutput.GetState().(states.TemperaturePortState).TStag
}

func (node *compressorNode) pStagIn() float64 {
	return node.pressureInput.GetState().(states.PressurePortState).PStag
}

func (node *compressorNode) pStagOut() float64 {
	return node.pressureOutput.GetState().(states.PressurePortState).PStag
}

func (node *compressorNode) gas() gases.Gas {
	return node.gasInput.GetState().(states.GasPortState).Gas
}
