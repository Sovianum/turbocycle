package constructive

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/Sovianum/turbocycle/impl/nodes"
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
}

type regeneratorNode struct {
	ports     core.PortsType
	sigma     float64
	precision float64
	mode      string
}

func NewRegeneratorNode(sigma, precision float64, mode string) RegeneratorNode {
	var result = &regeneratorNode{
		ports:     make(core.PortsType),
		sigma:     sigma,
		precision: precision,
		mode:      mode,
	}

	result.ports[nodes.HotGasInput] = core.NewPort()
	result.ports[nodes.HotGasInput].SetInnerNode(result)
	result.ports[nodes.HotGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.ColdGasInput] = core.NewPort()
	result.ports[nodes.ColdGasInput].SetInnerNode(result)
	result.ports[nodes.ColdGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.HotGasOutput] = core.NewPort()
	result.ports[nodes.HotGasOutput].SetInnerNode(result)
	result.ports[nodes.HotGasOutput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.ColdGasOutput] = core.NewPort()
	result.ports[nodes.ColdGasOutput].SetInnerNode(result)
	result.ports[nodes.ColdGasOutput].SetState(states.StandardAtmosphereState())

	return result
}

func (node *regeneratorNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		HotInputState   core.PortState `json:"hot_input_state"`
		ColdInputState  core.PortState `json:"cold_input_state"`
		HotOutputState  core.PortState `json:"hot_output_state"`
		ColdOutputState core.PortState `json:"cold_output_state"`
		Sigma           float64        `json:"sigma"`
	}{
		HotInputState:   node.hotInput().GetState(),
		ColdInputState:  node.coldInput().GetState(),
		HotOutputState:  node.hotOutput().GetState(),
		ColdOutputState: node.coldOutput().GetState(),
		Sigma:           node.sigma,
	})
}

func (node *regeneratorNode) ContextDefined() bool {
	return true
}

func (node *regeneratorNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *regeneratorNode) ColdInput() core.Port {
	return node.coldInput()
}

func (node *regeneratorNode) ColdOutput() core.Port {
	return node.coldOutput()
}

func (node *regeneratorNode) HotInput() core.Port {
	return node.hotInput()
}

func (node *regeneratorNode) HotOutput() core.Port {
	return node.hotOutput()
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
		return errors.New(fmt.Sprintf("Invalid Regenerator node state: %s", node.mode))
	}

	coldGasState.TStag, hotGasState.TStag = node.getNewTOut(coldGasState.TStag, hotGasState.TStag, iterFunc)
	node.hotOutput().SetState(hotGasState)
	node.coldOutput().SetState(coldGasState)

	switch node.mode {
	case SigmaByColdSide:
		var hotInputState = node.hotInput().GetState().(states.ComplexGasPortState)
		hotInputState.PStag = node.hotOutput().GetState().(states.ComplexGasPortState).PStag
		node.hotInput().SetState(hotInputState)
	case SigmaByHotSide:
		var coldInputState = node.coldInput().GetState().(states.ComplexGasPortState)
		coldInputState.PStag = node.coldOutput().GetState().(states.ComplexGasPortState).PStag
		node.coldInput().SetState(coldInputState)
	}

	return nil
}

func (node *regeneratorNode) GetRequirePortTags() ([]string, error) {
	switch node.mode {
	case SigmaByColdSide:
		return []string{nodes.ColdGasInput, nodes.HotGasInput, nodes.HotGasOutput}, nil
	case SigmaByHotSide:
		return []string{nodes.ColdGasInput, nodes.HotGasInput, nodes.ColdGasOutput}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid Regenerator node state: %s", node.mode))
	}
}

func (node *regeneratorNode) GetUpdatePortTags() ([]string, error) {
	switch node.mode {
	case SigmaByColdSide:
		return []string{nodes.ColdGasOutput, nodes.HotGasOutput, nodes.HotGasInput}, nil
	case SigmaByHotSide:
		return []string{nodes.ColdGasOutput, nodes.HotGasOutput, nodes.ColdGasInput}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid Regenerator node state: %s", node.mode))
	}
}

func (node *regeneratorNode) GetPortTags() []string {
	return []string{nodes.ColdGasInput, nodes.ColdGasOutput, nodes.HotGasInput, nodes.HotGasOutput}
}

func (node *regeneratorNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ColdGasInput:
		return node.coldInput(), nil
	case nodes.ColdGasOutput:
		return node.coldOutput(), nil
	case nodes.HotGasInput:
		return node.hotInput(), nil
	case nodes.HotGasOutput:
		return node.hotOutput(), nil
	default:
		return nil, errors.New(fmt.Sprintf(
			"port \"%s\" not found on regeneratorNode", tag,
		))
	}
}

func (node *regeneratorNode) getNewTOut(
	tStagColdOutCurr, tStagHotOutCurr float64,
	iterFunc func(float64, float64) (float64, float64),
) (tStagColdOutNew, tStagHotOutNew float64) {
	tStagColdOutNew, tStagHotOutNew = iterFunc(tStagColdOutCurr, tStagHotOutCurr)

	for !(common.Converged(tStagColdOutCurr, tStagColdOutNew, node.precision) && common.Converged(tStagHotOutCurr, tStagHotOutNew, node.precision)) {
		tStagColdOutCurr, tStagHotOutCurr = tStagColdOutNew, tStagHotOutNew
		tStagColdOutNew, tStagHotOutNew = iterFunc(tStagColdOutCurr, tStagHotOutCurr)
	}
	return
}

func (node *regeneratorNode) getNewTOutSigmaByColdSide(tStagColdOutCurr, tStagHotOutCurr float64) (tStagColdOut, tStagHotOut float64) {
	var coldGasState = node.coldInput().GetState().(states.ComplexGasPortState)
	var hotGasState = node.HotInput().GetState().(states.ComplexGasPortState)

	var hotHeatRate = hotGasState.MassRateRel * gases.CpMean(hotGasState.Gas, hotGasState.TStag, tStagHotOutCurr, nodes.DefaultN)
	var coldHeatRate = coldGasState.MassRateRel * gases.CpMean(coldGasState.Gas, coldGasState.TStag, tStagColdOutCurr, nodes.DefaultN)
	var heatRateFactor = hotHeatRate / coldHeatRate

	tStagColdOut = coldGasState.TStag + node.sigma*(hotGasState.TStag-coldGasState.TStag)
	tStagHotOut = node.sigma/heatRateFactor*coldGasState.TStag + (1-node.sigma/heatRateFactor)*hotGasState.TStag
	return
}

func (node *regeneratorNode) getNewTOutSigmaByHotSide(tStagColdOutCurr, tStagHotOutCurr float64) (tStagColdOut, tStagHotOut float64) {
	var coldGasState = node.coldInput().GetState().(states.ComplexGasPortState)
	var hotGasState = node.HotInput().GetState().(states.ComplexGasPortState)

	var hotHeatRate = hotGasState.MassRateRel * gases.CpMean(hotGasState.Gas, hotGasState.TStag, tStagHotOutCurr, nodes.DefaultN)
	var coldHeatRate = coldGasState.MassRateRel * gases.CpMean(coldGasState.Gas, coldGasState.TStag, tStagColdOutCurr, nodes.DefaultN)
	var heatRateFactor = hotHeatRate / coldHeatRate

	tStagHotOut = coldGasState.TStag + node.sigma*(hotGasState.TStag-coldGasState.TStag)
	tStagColdOut = coldGasState.TStag + heatRateFactor*node.sigma*(hotGasState.TStag-coldGasState.TStag)
	return
}

func (node *regeneratorNode) tStagHotIn() float64 {
	return node.hotInput().GetState().(states.ComplexGasPortState).TStag
}

func (node *regeneratorNode) tStagHotOut() float64 {
	return node.hotOutput().GetState().(states.ComplexGasPortState).TStag
}

func (node *regeneratorNode) tStagColdIn() float64 {
	return node.coldInput().GetState().(states.ComplexGasPortState).TStag
}

func (node *regeneratorNode) tStagColdOut() float64 {
	return node.coldOutput().GetState().(states.ComplexGasPortState).TStag
}

func (node *regeneratorNode) coldInput() core.Port {
	return node.ports[nodes.ColdGasInput]
}

func (node *regeneratorNode) coldOutput() core.Port {
	return node.ports[nodes.ColdGasOutput]
}

func (node *regeneratorNode) hotInput() core.Port {
	return node.ports[nodes.HotGasInput]
}

func (node *regeneratorNode) hotOutput() core.Port {
	return node.ports[nodes.HotGasOutput]
}
