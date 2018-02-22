package methodics

import (
	"testing"

	"github.com/stretchr/testify/suite"
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
	tc *turbineCharacteristic
}

func (s *TurbineCharacteristicTestSuite) SetupTest() {
	s.tc = NewTurbineCharacteristic(etaT0, piT0, massRateNormTurb0, lambdaU0, stageNum).(*turbineCharacteristic)
}

func (s *TurbineCharacteristicTestSuite) TestGetNormMassRateCharSamePoint() {
	c := s.tc.GetNormMassRateChar()
	mrn := c(lambdaU0, 1)
	s.InDelta(massRateNormTurb0, mrn, 1e-7)
}

func (s *TurbineCharacteristicTestSuite) TestGetNormEtaCharSamePoint() {
	c := s.tc.GetNormEtaChar()
	etaNorm := c(lambdaU0, 1)
	s.InDelta(etaT0, etaNorm, 1e-7)
}

func TestTurbineCharacteristicTestSuite(t *testing.T) {
	suite.Run(t, new(TurbineCharacteristicTestSuite))
}
