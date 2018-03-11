package methodics

import (
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
)

const (
	piC0          = 4
	etaStag0      = 0.86
	massRateNorm0 = 20
	etaC0         = 0.84
	rpmC0         = 10000
	p0            = 1e5
	t0            = 277
	precision     = 1e-4
)

type CompressorCharGenTestSuite struct {
	suite.Suite
	ccg    *compressorCharGen
	c      constructive.ParametricCompressorNode
	proto  constructive.CompressorNode
	protoC constructive.ParametricCompressorNode
}

func (s *CompressorCharGenTestSuite) SetupTest() {
	s.ccg = NewCompressorCharGen(
		piC0, etaStag0, massRateNorm0, 1e-6, 1, 100,
	).(*compressorCharGen)
	s.c = constructive.NewParametricCompressorNode(
		massRateNorm0, piC0, rpmC0, etaC0, t0, p0, 0.01,
		s.ccg.GetNormEtaChar(), s.ccg.GetNormRPMChar(),
	)
	s.c.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	s.c.TemperatureInput().SetState(states.NewTemperaturePortState(t0))
	s.c.PressureInput().SetState(states.NewPressurePortState(p0))
	s.c.MassRateInput().SetState(states.NewMassRatePortState(massRateNorm0))

	s.proto = constructive.NewCompressorNode(etaC0, piC0, precision)
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()),
			states.NewTemperaturePortState(t0),
			states.NewPressurePortState(p0),
			states.NewMassRatePortState(1),
		},
		[]graph.Port{
			s.proto.GasInput(), s.proto.TemperatureInput(),
			s.proto.PressureInput(), s.proto.MassRateInput(),
		},
	)
	s.Require().Nil(s.proto.Process())

	s.protoC = constructive.NewParametricCompressorNodeFromProto(
		s.proto, s.ccg.GetNormEtaChar(), s.ccg.GetNormRPMChar(), rpmC0, massRateNorm0, precision,
	)
	s.Require().Nil(s.protoC.Process())
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

func (s *CompressorCharGenTestSuite) TestConsistency() {
	s.InDelta(s.proto.TStagIn(), s.protoC.TStagIn(), 1e-4)
	s.InDelta(s.proto.PStagIn(), s.protoC.PStagIn(), 1e-4)
	s.InDelta(s.proto.TStagOut(), s.protoC.TStagOut(), 1e-4)
	s.InDelta(s.proto.PStagOut(), s.protoC.PStagOut(), 1e-4)
	s.InDelta(s.proto.PiStag(), s.protoC.PiStag(), 1e-4)
	s.InDelta(s.proto.Eta(), s.protoC.Eta(), 1e-9)
	s.InDelta(
		s.proto.PowerOutput().GetState().Value().(float64),
		s.protoC.PowerOutput().GetState().Value().(float64), 1e-1,
	)
	s.InDelta(massRateNorm0, s.protoC.MassRateInput().GetState().Value().(float64), 1e-4)
}

func TestCompressorCharGenTestSuite(t *testing.T) {
	suite.Run(t, new(CompressorCharGenTestSuite))
}
