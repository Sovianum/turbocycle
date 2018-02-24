package methodics

import (
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

const (
	piC0          = 4
	etaStag0      = 0.86
	massRateNorm0 = 20
)

type CompressorCharGenTestSuite struct {
	suite.Suite
	ccg *compressorCharGen
	c   constructive.ParametricCompressorNode
}

func (s *CompressorCharGenTestSuite) SetupTest() {
	s.ccg = NewCompressorCharGen(
		piC0, etaStag0, massRateNorm0, 1e-6, 1, 100,
	).(*compressorCharGen)
	s.c = constructive.NewParametricCompressorNode(
		massRateNorm0, piC0, 10000, 0.84, 300, 1e5, 0.01,
		s.ccg.GetNormEtaChar(), s.ccg.GetNormRPMChar(),
	)
	s.c.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	s.c.TemperatureInput().SetState(states.NewTemperaturePortState(300))
	s.c.PressureInput().SetState(states.NewPressurePortState(1e5))
	s.c.MassRateInput().SetState(states.NewMassRatePortState(massRateNorm0))
}

func (s *CompressorCharGenTestSuite) TestEtaStag() {
	etaStag := s.ccg.etaStag(math.Pi/4, 1)
	s.InDelta(s.ccg.etaStag0, etaStag, 1e-7)
}

func (s *CompressorCharGenTestSuite) TestPiCOpt() {
	pi := s.ccg.piCOpt(1)
	s.InDelta(piC0, pi, 1e-7)
}

func (s *CompressorCharGenTestSuite) TestMassRateNormOpt() {
	massRateNorm := s.ccg.massRateNorm(math.Pi/4, 1)
	s.InDelta(massRateNorm0, massRateNorm, 1e-7)
}

func (s *CompressorCharGenTestSuite) TestSolveEtaOrigin() {
	etaChar := s.ccg.GetNormEtaChar()
	etaNorm := etaChar(1, 1)
	s.InDelta(1, etaNorm, 1e-7)
}

func (s *CompressorCharGenTestSuite) TestSolveRPMOrigin() {
	rpmChar := s.ccg.GetNormRPMChar()
	rpmNorm := rpmChar(1, 1)
	s.InDelta(1, rpmNorm, 1e-7)
}

func (s *CompressorCharGenTestSuite) TestSolveNear() {
	s.Require().Nil(s.c.Process())
	sg := newton.NewUniformNewtonSolverGen(1e-4, newton.NoLog)
	currRpm := s.c.RPM()
	currEta := s.c.Eta()

	coef := 0.999999999
	expectedRpm := currRpm * coef
	expectedEta := currEta * coef

	vs := variator.NewVariatorSolver(
		func() (*mat.VecDense, error) {
			s.c.Process()
			return mat.NewVecDense(2, []float64{
				s.c.RPM() - expectedRpm,
				s.c.Eta() - expectedEta,
			}), nil
		},
		[]variator.Variator{
			variator.FromCallables(s.c.NormPiStag, s.c.SetNormPiStag),
			variator.FromCallables(s.c.NormMassRate, s.c.SetNormMassRate),
		},
		sg,
	)
	_, err := vs.Solve(vs.GetInit(), 1e-7, 0.1, 1000)
	s.Require().Nil(err)
	s.InDelta(expectedRpm, s.c.RPM(), 1e-6)
	s.InDelta(expectedEta, s.c.Eta(), 1e-6)
}

func (s *CompressorCharGenTestSuite) TestSolveConvergence() {
	s.Require().Nil(s.c.Process())
	sg := newton.NewUniformNewtonSolverGen(1e-4, newton.NoLog)
	vs := variator.NewVariatorSolver(
		func() (*mat.VecDense, error) {
			s.c.Process()
			return mat.NewVecDense(2, []float64{
				s.c.RPM() - 8000,
				s.c.Eta() - etaStag0*0.8,
			}), nil
		},
		[]variator.Variator{
			variator.FromCallables(s.c.NormPiStag, s.c.SetNormPiStag),
			variator.FromCallables(s.c.NormMassRate, s.c.SetNormMassRate),
		},
		sg,
	)
	_, err := vs.Solve(vs.GetInit(), 1e-7, 0.1, 1000)
	s.Require().Nil(err)
	s.InDelta(8000, s.c.RPM(), 1e-6)
	s.InDelta(etaStag0*0.8, s.c.Eta(), 1e-6)
}

func TestCompressorCharGenTestSuite(t *testing.T) {
	suite.Run(t, new(CompressorCharGenTestSuite))
}
