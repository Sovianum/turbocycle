package constructive

import (
	"fmt"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

const (
	SigmaByHotSide  = "sigmaByHotSide"
	SigmaByColdSide = "sigmaByColdSide"
)

type RegeneratorNode interface {
	graph.Node
	ColdInput() graph.Port
	ColdOutput() graph.Port
	HotInput() graph.Port
	HotOutput() graph.Port
	Sigma() float64
}

func NewRegeneratorNode(sigma, precision float64, mode string) RegeneratorNode {
	var result = &regeneratorNode{
		sigma:     sigma,
		precision: precision,
		mode:      mode,
	}

	result.hotGasInput = graph.NewAttachedPort(result)
	result.coldGasInput = graph.NewAttachedPort(result)
	result.hotGasOutput = graph.NewAttachedPort(result)
	result.coldGasOutput = graph.NewAttachedPort(result)

	return result
}

type regeneratorNode struct {
	graph.BaseNode

	hotGasInput   graph.Port
	coldGasInput  graph.Port
	hotGasOutput  graph.Port
	coldGasOutput graph.Port

	sigma     float64
	precision float64
	mode      string
}

func (node *regeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Regenerator")
}

func (node *regeneratorNode) GetPorts() []graph.Port {
	return []graph.Port{node.hotGasInput, node.coldGasInput, node.hotGasOutput, node.coldGasOutput}
}

func (node *regeneratorNode) GetRequirePorts() []graph.Port {
	return []graph.Port{node.hotGasInput, node.coldGasInput}
}

func (node *regeneratorNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{node.hotGasOutput, node.coldGasOutput}
}

func (node *regeneratorNode) ColdInput() graph.Port {
	return node.coldGasInput
}

func (node *regeneratorNode) ColdOutput() graph.Port {
	return node.coldGasOutput
}

func (node *regeneratorNode) HotInput() graph.Port {
	return node.hotGasInput
}

func (node *regeneratorNode) HotOutput() graph.Port {
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
