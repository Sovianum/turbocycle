package methodics

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
)

const (
	etaT0             = 0.92
	piT0              = 1.4
	massRateNormTurb0 = 10
	lambdaU0          = 0.5
	stageNum          = 1
)

type TurbineCharacteristicTestSuite struct {
	suite.Suite
	tc  *ciamTurbineCharacteristic
	ccg *compressorCharGen
	t   constructive.ParametricTurbineNode
	c   constructive.ParametricCompressorNode
}

func (s *TurbineCharacteristicTestSuite) SetupTest() {
	s.tc = NewCIAMTurbineCharacteristic(etaT0, piT0, lambdaU0, stageNum).(*ciamTurbineCharacteristic)
	s.t = constructive.NewParametricTurbineNode(
		massRateNormTurb0,
		piT0, etaT0, 500, 1e6, 0.5, 0.01,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		s.tc.GetNormMassRateChar(),
		s.tc.GetNormEtaChar(),
	)
	s.t.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	s.t.TemperatureInput().SetState(states.NewTemperaturePortState(1200))
	s.t.PressureInput().SetState(states.NewPressurePortState(1e6))
	s.t.MassRateInput().SetState(states.NewMassRatePortState(massRateNormTurb0))
	s.t.RPMInput().SetState(states.NewRPMPortState(10000))

	s.ccg = NewCompressorCharGen(
		piC0, etaStag0, massRateNorm0, 1e-6, 1, 100,
	).(*compressorCharGen)
	s.c = constructive.NewParametricCompressorNode(
		massRateNorm0, piC0, 10000, 0.84, 300, 1e5, 0.01,
		s.ccg.GetNormEtaChar(),
		s.ccg.GetNormRPMChar(),
	)
	s.c.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	s.c.TemperatureInput().SetState(states.NewTemperaturePortState(300))
	s.c.PressureInput().SetState(states.NewPressurePortState(1e5))
	s.c.MassRateInput().SetState(states.NewMassRatePortState(massRateNorm0))

	graph.Link(s.t.RPMInput(), s.c.RPMOutput())
}

func (s *TurbineCharacteristicTestSuite) TestGetNormMassRateCharSamePoint() {
	c := s.tc.GetNormMassRateChar()
	mrn := c(lambdaU0, 1)
	s.InDelta(1, mrn, 1e-7)
}

func (s *TurbineCharacteristicTestSuite) TestGetNormEtaCharSamePoint() {
	c := s.tc.GetNormEtaChar()
	etaNorm := c(lambdaU0, 1)
	s.InDelta(etaT0, etaNorm, 1e-7)
}

func (s *TurbineCharacteristicTestSuite) TestTrends() {
	tc := NewCIAMTurbineCharacteristic(etaT0, piT0, lambdaU0, stageNum)
	normEtaFunc := tc.GetNormEtaChar()
	normMassRateFunc := tc.GetNormMassRateChar()

	s.True(normEtaFunc(lambdaU0*1, 0.99) < 1, "%f", normEtaFunc(lambdaU0, 0.99))
	s.True(normMassRateFunc(lambdaU0, 0.99) < 1, "%f", normMassRateFunc(lambdaU0, 0.99))

	s.True(normEtaFunc(lambdaU0*1, 1.01) < 1, "%f", normEtaFunc(lambdaU0, 1.01))
	s.True(normMassRateFunc(lambdaU0*1, 1.01) >= 1, "%f", normMassRateFunc(lambdaU0, 1.01))

	s.True(normEtaFunc(lambdaU0*0.99, 1) < 1, "%f", normEtaFunc(lambdaU0, 0.99))

	s.True(normEtaFunc(lambdaU0*1.01, 1) < 1, "%f", normEtaFunc(lambdaU0, 1.01))
}

func TestTurbineCharacteristicTestSuite(t *testing.T) {
	suite.Run(t, new(TurbineCharacteristicTestSuite))
}
