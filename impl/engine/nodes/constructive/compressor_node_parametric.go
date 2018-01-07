package constructive

import (
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type CompressorCharFunc func(normMassRate, normPiStag float64) float64

type ParametricCompressorNode interface {
	graph.Node

	nodes.TemperatureChannel
	nodes.PressureChannel
	nodes.GasChannel
	nodes.PowerSource
	nodes.RPMSource
	nodes.MassRateChannel

	nodes.PressureIn
	nodes.PressureOut
	nodes.TemperatureIn
	nodes.TemperatureOut

	LSpecific() float64
	Eta() float64

	PiStag() float64
	SetPiStag(piStag float64)
	NormPiStag() float64
	SetNormPiStag(normPiStag float64)

	MassRate() float64
	NormMassRate() float64
	SetNormMassRate(normMassRate float64)

	NormalizedRPM() float64
	RPM() float64
}

func NewParametricCompressorNode(
	massRate0, piStag0, rpm0, eta0, t0, p0, precision float64,
	normEtaCharacteristic, normRpmCharacteristic CompressorCharFunc,
) ParametricCompressorNode {
	var result = &parametricCompressorNode{
		// normalized parameters are equal to the ones  at nominal mode
		normPiStag:   1,
		normMassRate: 1,

		t0: t0,
		p0: p0,

		massRate0: massRate0,
		piStag0:   piStag0,
		rpm0:      rpm0,
		eta0:      eta0,

		precision: precision,

		normEtaCharacteristic: normEtaCharacteristic,
		normRpmCharacteristic: normRpmCharacteristic,
	}

	result.baseCompressor = newBaseCompressor(result, precision)
	result.rpmOutput = graph.NewAttachedPort(result)
	result.massRateInput = graph.NewAttachedPort(result)

	return result
}

type parametricCompressorNode struct {
	*baseCompressor

	rpmOutput     graph.Port
	massRateInput graph.Port

	normEtaCharacteristic CompressorCharFunc
	normRpmCharacteristic CompressorCharFunc

	normPiStag   float64
	normMassRate float64

	t0 float64
	p0 float64

	massRate0 float64
	rpm0      float64
	eta0      float64
	piStag0   float64

	precision float64
}

func (node *parametricCompressorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ParametricCompressor")
}

func (node *parametricCompressorNode) Process() error {
	if node.piStag() <= 1 {
		return fmt.Errorf("invalid piStag = %f", node.piStag)
	}

	var etaAd = node.normEtaCharacteristic(node.normMassRate, node.piStag()) * node.eta0
	var tStagOut, err = node.getTStagOut(node.tStagIn(), node.piStag(), etaAd)
	if err != nil {
		return err
	}
	var pStagOut = node.pStagIn() * node.piStag()

	var massRate = node.massRate()
	graph.SetAll(
		[]graph.PortState{
			states.NewMassRatePortState(massRate),
			node.gasInput.GetState(),
			states.NewTemperaturePortState(tStagOut),
			states.NewPressurePortState(pStagOut),
			states.NewMassRatePortState(massRate),
		},
		[]graph.Port{
			node.massRateInput,
			node.gasOutput,
			node.temperatureOutput,
			node.pressureOutput,
			node.massRateOutput,
		},
	)
	node.powerOutput.SetState(
		states.NewPowerPortState(-node.lSpecific()),
	)
	node.rpmOutput.SetState(
		states.NewRPMPortState(node.rpm()),
	)

	return nil
}

func (node *parametricCompressorNode) GetPorts() []graph.Port {
	return append(node.baseCompressor.GetPorts(), node.rpmOutput, node.massRateInput)
}

func (node *parametricCompressorNode) GetUpdatePorts() []graph.Port {
	return append(node.baseCompressor.GetUpdatePorts(), node.rpmOutput, node.massRateInput)
}

func (node *parametricCompressorNode) RPMOutput() graph.Port {
	return node.rpmOutput
}

func (node *parametricCompressorNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *parametricCompressorNode) PiStag() float64 {
	return node.piStag()
}

func (node *parametricCompressorNode) SetPiStag(piStag float64) {
	node.normPiStag = piStag / node.piStag0
}

func (node *parametricCompressorNode) NormPiStag() float64 {
	return node.normPiStag
}

func (node *parametricCompressorNode) SetNormPiStag(normPiStag float64) {
	node.normPiStag = normPiStag
}

func (node *parametricCompressorNode) MassRate() float64 {
	return node.massRate()
}

func (node *parametricCompressorNode) NormMassRate() float64 {
	return node.normMassRate
}

func (node *parametricCompressorNode) SetNormMassRate(normMassRate float64) {
	node.normMassRate = normMassRate
}

func (node *parametricCompressorNode) Eta() float64 {
	return node.normEtaCharacteristic(node.normMassRate, node.normPiStag) * node.eta0
}

func (node *parametricCompressorNode) RPM() float64 {
	return node.rpm()
}

func (node *parametricCompressorNode) NormalizedRPM() float64 {
	return node.normRpmCharacteristic(node.normMassRate, node.normPiStag)
}

func (node *parametricCompressorNode) rpm() float64 {
	var tFactor = math.Sqrt(node.t0 / node.tStagIn())
	return node.normRpmCharacteristic(node.normMassRate, node.normPiStag) / tFactor * node.rpm0
}

func (node *parametricCompressorNode) massRate() float64 {
	var tFactor = math.Sqrt(node.t0 / node.tStagIn())
	var pFactor = node.pStagIn() / node.p0
	return node.normMassRate / (tFactor * pFactor) * node.massRate0
}

func (node *parametricCompressorNode) piStag() float64 {
	return node.normPiStag * node.piStag0
}
