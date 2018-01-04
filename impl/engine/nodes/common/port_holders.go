package common

import "github.com/Sovianum/turbocycle/core/graph"

func NewPressureChannelHolder(node graph.Node) *PressureChannelHolder {
	return &PressureChannelHolder{
		PressureSinkHolder:   NewPressureSinkHolder(node),
		PressureSourceHolder: NewPressureSourceHolder(node),
	}
}

type PressureChannelHolder struct {
	*PressureSinkHolder
	*PressureSourceHolder
}

func NewPressureSourceHolder(node graph.Node) *PressureSourceHolder {
	var h = &PressureSourceHolder{}
	h.pressureOutput = graph.NewAttachedPort(node)
	return h
}

type PressureSourceHolder struct {
	pressureOutput graph.Port
}

func (h *PressureSourceHolder) PressureOutput() graph.Port {
	return h.pressureOutput
}

func NewPressureSinkHolder(node graph.Node) *PressureSinkHolder {
	var h = &PressureSinkHolder{}
	h.pressureInput = graph.NewAttachedPort(node)
	return h
}

type PressureSinkHolder struct {
	pressureInput graph.Port
}

func (h *PressureSinkHolder) PressureInput() graph.Port {
	return h.pressureInput
}
