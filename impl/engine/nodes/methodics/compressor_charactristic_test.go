package methodics

import (
	"math"
	"testing"

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
}

func (s *CompressorCharGenTestSuite) SetupTest() {
	s.ccg = NewCompressorCharGen(
		piC0, etaStag0, massRateNorm0, 1e-6, 1, 100,
	).(*compressorCharGen)
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

func TestCompressorCharGenTestSuite(t *testing.T) {
	suite.Run(t, new(CompressorCharGenTestSuite))
}
