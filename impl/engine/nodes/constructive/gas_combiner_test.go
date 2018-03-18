package constructive

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
)

type GasCombinerTestSuite struct {
	suite.Suite
	gas      gases.Gas
	cp       float64
	combiner GasCombiner

	mainT  float64
	mainP  float64
	mainMR float64

	extraT  float64
	extraP  float64
	extraMR float64
}

func (s *GasCombinerTestSuite) SetupTest() {
	s.cp = 1000

	s.mainT = 1000
	s.mainP = 1e6
	s.mainMR = 50

	s.extraT = 500
	s.extraP = 1e7
	s.extraMR = 10

	s.gas = &gases.TestGas{CpVal: s.cp}
	s.combiner = NewGasCombiner(1e-6, 1, 10)

	mainInput := s.combiner.MainInput()
	extraInput := s.combiner.ExtraInput()

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(s.gas),
			states.NewTemperaturePortState(s.mainT),
			states.NewPressurePortState(s.mainP),
			states.NewMassRatePortState(s.mainMR),

			states.NewGasPortState(s.gas),
			states.NewTemperaturePortState(s.extraT),
			states.NewPressurePortState(s.extraP),
			states.NewMassRatePortState(s.extraMR),
		},
		[]graph.Port{
			mainInput.GasInput(), mainInput.TemperatureInput(), mainInput.PressureInput(), mainInput.MassRateInput(),
			extraInput.GasInput(), extraInput.TemperatureInput(), extraInput.PressureInput(), extraInput.MassRateInput(),
		},
	)
}

func (s *GasCombinerTestSuite) TestProcess() {
	s.Require().Nil(s.combiner.Process())
	s.InDelta(s.cp, s.combiner.Output().GasOutput().GetState().Value().(gases.Gas).Cp(0), 1e-9)
	s.InDelta(s.mainMR+s.extraMR, s.combiner.Output().MassRateOutput().GetState().Value().(float64), 1e-9)

	mainFraction := s.mainMR / (s.mainMR + s.extraMR)
	extraFraction := s.extraMR / (s.mainMR + s.extraMR)
	tOut := mainFraction*s.mainT + extraFraction*s.extraT
	s.InDelta(tOut, s.combiner.Output().TemperatureOutput().GetState().Value().(float64), 1e-9)

	s.InDelta(s.mainP, s.combiner.Output().PressureOutput().GetState().Value().(float64), 1e-9)
}

func TestGasCombinerTestSuite(t *testing.T) {
	suite.Run(t, new(GasCombinerTestSuite))
}
