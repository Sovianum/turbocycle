package methodics

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

const (
	etaT0             = 0.92
	piT0              = 3
	massRateNormTurb0 = 10
	lambdaU0          = 0.5
	stageNum          = 1
)

type TurbineCharacteristicTestSuite struct {
	suite.Suite
	tc  *turbineCharacteristic
	ccg *compressorCharGen
	t   constructive.ParametricTurbineNode
	c   constructive.ParametricCompressorNode
}

func (s *TurbineCharacteristicTestSuite) SetupTest() {
	s.tc = NewTurbineCharacteristic(etaT0, piT0, lambdaU0, stageNum).(*turbineCharacteristic)
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

func (s *TurbineCharacteristicTestSuite) TestSolveConvergence() {
	s.Require().Nil(s.t.Process())
	sg := newton.NewUniformNewtonSolverGen(1e-4, newton.NoLog)

	startPower := s.t.PowerOutput().GetState().Value().(float64)
	expectedPower := startPower * 0.8
	vs := variator.NewVariatorSolver(
		func() (*mat.VecDense, error) {
			s.t.Process()
			return mat.NewVecDense(1, []float64{
				s.t.PowerOutput().GetState().Value().(float64) - expectedPower,
			}), nil
		},
		[]variator.Variator{variator.FromCallables(s.t.NormPiT, s.t.SetNormPiT)},
		sg,
	)
	_, err := vs.Solve(vs.GetInit(), 1e-6, 0.01, 10000)
	s.Require().Nil(err)

	s.InDelta(expectedPower, s.t.PowerOutput().GetState().Value().(float64), 1e-6)
}

func (s *TurbineCharacteristicTestSuite) TestComplexConverge() {
	s.Require().Nil(s.t.Process())
	s.Require().Nil(s.c.Process())
	sg := newton.NewUniformNewtonSolverGen(1e-4, newton.NoLog)

	expectedPOut := s.t.PStagOut() * 1.1

	vs := variator.NewVariatorSolver(
		func() (*mat.VecDense, error) {
			s.t.Process()
			s.c.Process()
			return mat.NewVecDense(3, []float64{
				s.c.MassRate() - s.t.MassRateInput().GetState().Value().(float64),
				s.c.MassRate()*s.c.PowerOutput().GetState().Value().(float64) +
					s.t.MassRateInput().GetState().Value().(float64)*
						s.t.PowerOutput().GetState().Value().(float64),
				s.t.PStagOut() - expectedPOut,
			}), nil
		},
		[]variator.Variator{
			variator.FromCallables(s.t.NormPiT, s.t.SetNormPiT),
			variator.FromCallables(s.c.NormPiStag, s.c.SetNormPiStag),
			variator.FromCallables(s.c.NormMassRate, s.c.SetNormMassRate),
		},
		sg,
	)
	_, err := vs.Solve(vs.GetInit(), 1e-6, 0.2, 10000)
	s.Require().Nil(err)

	s.InDelta(s.t.PStagOut(), expectedPOut, 1e-6)
}

func TestTurbineCharacteristicTestSuite(t *testing.T) {
	suite.Run(t, new(TurbineCharacteristicTestSuite))
}
