package methodics

import (
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
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

func (s *CompressorCharGenTestSuite) TestTrends() {
	s.Require().Nil(s.c.Process())

	rpm0 := s.c.RPM()
	eta0 := s.c.Eta()
	pi0 := s.c.PiStag()
	normMassRate0 := s.c.NormMassRate()

	init := func(piFrac, nmrFrac float64) error {
		s.c.SetPiStag(piFrac * pi0)
		s.c.SetNormMassRate(nmrFrac * normMassRate0)
		return s.c.Process()
	}

	s.Require().Nil(init(1.01, 1))
	s.True(s.c.RPM() > rpm0, "%f %f", s.c.RPM(), rpm0)
	s.True(s.c.Eta() < eta0, "%f %f", s.c.Eta(), eta0)

	s.Require().Nil(init(0.99, 1))
	s.True(s.c.RPM() < rpm0, "%f %f", s.c.RPM(), rpm0)
	s.True(s.c.Eta() < eta0, "%f %f", s.c.Eta(), eta0)

	s.Require().Nil(init(1, 1.01))
	s.True(s.c.RPM() > rpm0, "%f %f", s.c.RPM(), rpm0)
	s.True(s.c.Eta() < eta0, "%f %f", s.c.Eta(), eta0)

	s.Require().Nil(init(1, 0.99))
	s.True(s.c.RPM() < rpm0, "%f %f", s.c.RPM(), rpm0)
	s.True(s.c.Eta() < eta0, "%f %f", s.c.Eta(), eta0)
}

func TestCompressorCharGenTestSuite(t *testing.T) {
	suite.Run(t, new(CompressorCharGenTestSuite))
}
