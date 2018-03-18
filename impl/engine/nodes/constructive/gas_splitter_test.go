package constructive

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
)

func TestGasSplitter_Process(t *testing.T) {
	temp := 300.
	p := 1e5
	mr := 10.
	factor := 0.1

	gs := NewGasSplitter(factor)
	input := gs.Input()
	mOut := gs.MainOutput()
	eOut := gs.ExtraOutput()

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()),
			states.NewTemperaturePortState(temp),
			states.NewPressurePortState(p),
			states.NewMassRatePortState(mr),
		},
		[]graph.Port{
			input.GasInput(), input.TemperatureInput(),
			input.PressureInput(), input.MassRateInput(),
		},
	)

	assert.Nil(t, gs.Process())

	assert.InDelta(t, temp, mOut.TemperatureOutput().GetState().Value().(float64), 1e-9)
	assert.InDelta(t, p, mOut.PressureOutput().GetState().Value().(float64), 1e-9)
	assert.InDelta(t, mr*(1-factor), mOut.MassRateOutput().GetState().Value().(float64), 1e-9)
	assert.NotNil(t, mOut.GasOutput().GetState().Value().(gases.Gas))

	assert.InDelta(t, temp, eOut.TemperatureOutput().GetState().Value().(float64), 1e-9)
	assert.InDelta(t, p, eOut.PressureOutput().GetState().Value().(float64), 1e-9)
	assert.InDelta(t, mr*factor, eOut.MassRateOutput().GetState().Value().(float64), 1e-9)
	assert.NotNil(t, eOut.GasOutput().GetState().Value().(gases.Gas))
}
