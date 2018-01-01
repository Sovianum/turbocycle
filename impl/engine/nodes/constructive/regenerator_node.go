package constructive

import (
	"fmt"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

const (
	SigmaByHotSide  = "sigmaByHotSide"
	SigmaByColdSide = "sigmaByColdSide"
)

type RegeneratorNode interface {
	core.Node
	ColdInput() core.Port
	ColdOutput() core.Port
	HotInput() core.Port
	HotOutput() core.Port
	Sigma() float64
}

func NewRegeneratorNode(sigma, precision float64, mode string) RegeneratorNode {
	var result = &regeneratorNode{
		sigma:     sigma,
		precision: precision,
		mode:      mode,
	}

	result.hotGasInput = core.NewAttachedPort(result)
	result.coldGasInput = core.NewAttachedPort(result)
	result.hotGasOutput = core.NewAttachedPort(result)
	result.coldGasOutput = core.NewAttachedPort(result)

	return result
}

type regeneratorNode struct {
	core.BaseNode

	hotGasInput   core.Port
	coldGasInput  core.Port
	hotGasOutput  core.Port
	coldGasOutput core.Port

	sigma     float64
	precision float64
	mode      string
}

func (node *regeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Regenerator")
}

func (node *regeneratorNode) GetPorts() []core.Port {
	return []core.Port{node.hotGasInput, node.coldGasInput, node.hotGasOutput, node.coldGasOutput}
}

func (node *regeneratorNode) GetRequirePorts() []core.Port {
	return []core.Port{node.hotGasInput, node.coldGasInput}
}

func (node *regeneratorNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.hotGasOutput, node.coldGasOutput}
}

func (node *regeneratorNode) ColdInput() core.Port {
	return node.coldGasInput
}

func (node *regeneratorNode) ColdOutput() core.Port {
	return node.coldGasOutput
}

func (node *regeneratorNode) HotInput() core.Port {
	return node.hotGasInput
}

func (node *regeneratorNode) HotOutput() core.Port {
	return node.hotGasOutput
}

func (node *regeneratorNode) Sigma() float64 {
	return node.sigma
}

func (node *regeneratorNode) Process() error {
	var hotGasState = node.HotInput().GetState().(states.ComplexGasPortState)
	var coldGasState = node.ColdInput().GetState().(states.ComplexGasPortState)

	var iterFunc func(float64, float64) (float64, float64)
	switch node.mode {
	case SigmaByColdSide:
		iterFunc = node.getNewTOutSigmaByColdSide
	case SigmaByHotSide:
		iterFunc = node.getNewTOutSigmaByHotSide
	default:
		return fmt.Errorf("invalid Regenerator node state: %s", node.mode)
	}

	coldGasState.TStag, hotGasState.TStag = node.getNewTOut(coldGasState.TStag, hotGasState.TStag, iterFunc)
	node.hotGasOutput.SetState(hotGasState)
	node.coldGasOutput.SetState(coldGasState)

	switch node.mode {
	case SigmaByColdSide:
		var hotInputState = node.hotGasInput.GetState().(states.ComplexGasPortState)
		hotInputState.PStag = node.hotGasOutput.GetState().(states.ComplexGasPortState).PStag
		node.hotGasInput.SetState(hotInputState)
	case SigmaByHotSide:
		var coldInputState = node.coldGasInput.GetState().(states.ComplexGasPortState)
		coldInputState.PStag = node.coldGasOutput.GetState().(states.ComplexGasPortState).PStag
		node.coldGasInput.SetState(coldInputState)
	}

	return nil
}

func (node *regeneratorNode) getNewTOut(
	tStagColdOutCurr, tStagHotOutCurr float64,
	iterFunc func(float64, float64) (float64, float64),
) (tStagColdOutNew, tStagHotOutNew float64) {
	tStagColdOutNew, tStagHotOutNew = iterFunc(tStagColdOutCurr, tStagHotOutCurr)

	for !(common.Converged(tStagColdOutCurr, tStagColdOutNew, node.precision) &&
		common.Converged(tStagHotOutCurr, tStagHotOutNew, node.precision)) {
		tStagColdOutCurr, tStagHotOutCurr = tStagColdOutNew, tStagHotOutNew
		tStagColdOutNew, tStagHotOutNew = iterFunc(tStagColdOutCurr, tStagHotOutCurr)
	}
	return
}

func (node *regeneratorNode) getNewTOutSigmaByColdSide(tStagColdOutCurr, tStagHotOutCurr float64) (tStagColdOut, tStagHotOut float64) {
	var coldGasState = node.coldGasInput.GetState().(states.ComplexGasPortState)
	var hotGasState = node.hotGasInput.GetState().(states.ComplexGasPortState)

	var hotHeatRate = hotGasState.MassRateRel * gases.CpMean(hotGasState.Gas, hotGasState.TStag, tStagHotOutCurr, nodes.DefaultN)
	var coldHeatRate = coldGasState.MassRateRel * gases.CpMean(coldGasState.Gas, coldGasState.TStag, tStagColdOutCurr, nodes.DefaultN)
	var heatRateFactor = hotHeatRate / coldHeatRate

	tStagColdOut = coldGasState.TStag + node.sigma*(hotGasState.TStag-coldGasState.TStag)
	tStagHotOut = node.sigma/heatRateFactor*coldGasState.TStag + (1-node.sigma/heatRateFactor)*hotGasState.TStag
	return
}

func (node *regeneratorNode) getNewTOutSigmaByHotSide(tStagColdOutCurr, tStagHotOutCurr float64) (tStagColdOut, tStagHotOut float64) {
	var coldGasState = node.coldGasInput.GetState().(states.ComplexGasPortState)
	var hotGasState = node.hotGasInput.GetState().(states.ComplexGasPortState)

	var hotHeatRate = hotGasState.MassRateRel * gases.CpMean(hotGasState.Gas, hotGasState.TStag, tStagHotOutCurr, nodes.DefaultN)
	var coldHeatRate = coldGasState.MassRateRel * gases.CpMean(coldGasState.Gas, coldGasState.TStag, tStagColdOutCurr, nodes.DefaultN)
	var heatRateFactor = hotHeatRate / coldHeatRate

	tStagHotOut = coldGasState.TStag + node.sigma*(hotGasState.TStag-coldGasState.TStag)
	tStagColdOut = coldGasState.TStag + heatRateFactor*node.sigma*(hotGasState.TStag-coldGasState.TStag)
	return
}

func (node *regeneratorNode) tStagHotIn() float64 {
	return node.hotGasInput.GetState().(states.ComplexGasPortState).TStag
}

func (node *regeneratorNode) tStagHotOut() float64 {
	return node.hotGasOutput.GetState().(states.ComplexGasPortState).TStag
}

func (node *regeneratorNode) tStagColdIn() float64 {
	return node.coldGasInput.GetState().(states.ComplexGasPortState).TStag
}

func (node *regeneratorNode) tStagColdOut() float64 {
	return node.coldGasOutput.GetState().(states.ComplexGasPortState).TStag
}
