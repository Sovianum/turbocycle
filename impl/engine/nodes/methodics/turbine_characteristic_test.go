package methodics

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
)

const (
	etaT0             = 0.92
	piT0              = 1.4
	massRateNormTurb0 = 10
	lambdaU0          = 0.5
	stageNum          = 1
	alphaT            = 3
	tT0               = 500
	pT0               = 1e6
	diam              = 0.5
	rpmT0             = 10000
)

type TurbineCharacteristicTestSuite struct {
	suite.Suite
	tc           TurbineCharacteristic
	ccg          *compressorCharGen
	t            constructive.ParametricTurbineNode
	c            constructive.ParametricCompressorNode
	blockedProto constructive.BlockedTurbineNode
	freeProto    constructive.FreeTurbineNode
}

func (s *TurbineCharacteristicTestSuite) SetupTest() {
	s.tc = NewKazandjanTurbineCharacteristic()
	s.t = constructive.NewParametricTurbineNode(
		massRateNormTurb0,
		piT0, etaT0, tT0, pT0, diam, precision,
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
	s.t.TemperatureInput().SetState(states.NewTemperaturePortState(tT0))
	s.t.PressureInput().SetState(states.NewPressurePortState(pT0))
	s.t.MassRateInput().SetState(states.NewMassRatePortState(massRateNormTurb0))
	s.t.RPMInput().SetState(states.NewRPMPortState(rpmT0))

	s.ccg = NewCompressorCharGen(
		piC0, etaStag0, massRateNorm0, precision, 1, 100,
	).(*compressorCharGen)
	s.c = constructive.NewParametricCompressorNode(
		massRateNorm0, piC0, rpmT0, etaT0, tT0, pT0, precision,
		s.ccg.GetNormEtaChar(),
		s.ccg.GetNormRPMChar(),
	)
	s.c.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	s.c.TemperatureInput().SetState(states.NewTemperaturePortState(tT0))
	s.c.PressureInput().SetState(states.NewPressurePortState(pT0))
	s.c.MassRateInput().SetState(states.NewMassRatePortState(massRateNorm0))

	graph.Link(s.t.RPMInput(), s.c.RPMOutput())

	s.blockedProto = constructive.NewBlockedTurbineNode(
		etaT0, 0.3, 0.001,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
	)
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(fuel.GetCH4().GetCombustionGas(gases.GetAir(), alphaT)),
			states.NewPressurePortState(pT0),
			states.NewTemperaturePortState(tT0),
			states.NewMassRatePortState(1),
			states.NewPowerPortState(1.5e5),
		},
		[]graph.Port{
			s.blockedProto.GasInput(),
			s.blockedProto.PressureInput(),
			s.blockedProto.TemperatureInput(),
			s.blockedProto.MassRateInput(),
			s.blockedProto.PowerInput(),
		},
	)
	s.Require().Nil(s.blockedProto.Process())

	s.freeProto = constructive.NewFreeTurbineNode(
		etaT0, 0.3, precision,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
	)
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(fuel.GetCH4().GetCombustionGas(gases.GetAir(), alphaT)),
			states.NewPressurePortState(pT0),
			states.NewTemperaturePortState(tT0),
			states.NewMassRatePortState(1),
			states.NewPressurePortState(pT0 / 2),
		},
		[]graph.Port{
			s.freeProto.GasInput(),
			s.freeProto.PressureInput(),
			s.freeProto.TemperatureInput(),
			s.freeProto.MassRateInput(),
			s.freeProto.PressureOutput(),
		},
	)
	s.Require().Nil(s.freeProto.Process())
}

func (s *TurbineCharacteristicTestSuite) TestGetNormMassRateCharSamePoint() {
	c := s.tc.GetNormMassRateChar()
	mrn := c(lambdaU0, 1)
	s.InDelta(1, mrn, 1e-7)
}

func (s *TurbineCharacteristicTestSuite) TestGetNormEtaCharSamePoint() {
	c := s.tc.GetNormEtaChar()
	etaNorm := c(lambdaU0, 1)
	s.InDelta(1, etaNorm, 1e-7)
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

func (s *TurbineCharacteristicTestSuite) TestBlockedTurbineConsistency() {
	t := constructive.NewParametricTurbineNodeFromProto(
		s.blockedProto, s.tc.GetNormMassRateChar(), s.tc.GetNormEtaChar(), massRateNorm0, 0.5, 1e-4,
	)
	t.RPMInput().SetState(states.NewRPMPortState(rpmT0))
	s.Require().Nil(t.Process())

	s.InDelta(s.blockedProto.PStagIn(), t.PStagIn(), 1e-5)
	s.InDelta(s.blockedProto.TStagIn(), t.TStagIn(), 1e-5)
	s.InDelta(s.blockedProto.PStagOut(), t.PStagOut(), 1e-5)
	s.InDelta(s.blockedProto.TStagOut(), t.TStagOut(), 1e-2)
	s.InDelta(s.blockedProto.Eta(), t.Eta(), 1e-6)
	s.InDelta(
		s.blockedProto.PowerOutput().GetState().Value().(float64),
		t.PowerOutput().GetState().Value().(float64), 1e-1,
	)
	s.InDelta(s.blockedProto.PiTStag(), t.PiTStag(), 1e-6)
}

func (s *TurbineCharacteristicTestSuite) TestFreeTurbineConsistency() {
	t := constructive.NewParametricTurbineNodeFromProto(
		s.freeProto, s.tc.GetNormMassRateChar(), s.tc.GetNormEtaChar(), massRateNorm0, 0.5, 1e-4,
	)
	t.RPMInput().SetState(states.NewRPMPortState(rpmT0))
	s.Require().Nil(t.Process())

	s.InDelta(s.freeProto.PStagIn(), t.PStagIn(), 1e-5)
	s.InDelta(s.freeProto.TStagIn(), t.TStagIn(), 1e-5)
	s.InDelta(s.freeProto.PStagOut(), t.PStagOut(), 1e-5)
	s.InDelta(s.freeProto.TStagOut(), t.TStagOut(), 1e-2)
	s.InDelta(s.freeProto.Eta(), t.Eta(), 1e-6)
	s.InDelta(
		s.freeProto.PowerOutput().GetState().Value().(float64),
		t.PowerOutput().GetState().Value().(float64), 1e-1,
	)
}

func TestTurbineCharacteristicTestSuite(t *testing.T) {
	suite.Run(t, new(TurbineCharacteristicTestSuite))
}
