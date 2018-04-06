package compressor

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"github.com/Sovianum/turbocycle/impl/stage/ditributions"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

const (
	expectedPi  = 4
	expectedEta = 0.82
)

type CompressorFittingTetSuite struct {
	suite.Suite
	node   StagedCompressorNode
	setter func(f common.Func1D)
}

func (suite *CompressorFittingTetSuite) SetupTest() {
	suite.node = getTestCompressor()
	suite.Require().NoError(suite.node.Process())
	suite.setter = func(f common.Func1D) { suite.node.SetHtLaw(common.FromDistribution(f)) }
}

func (suite CompressorFittingTetSuite) TestPiUnitShape() {
	fg := common.Scaler(ditributions.GetUnitConstant())
	eqSys := GetPiFitEqSys1D(suite.node, expectedPi, suite.setter, fg)

	solver, err := newton.NewUniformNewtonSolver(eqSys, 1e-5, newton.NoLog)
	suite.Require().NoError(err)

	_, err = solver.Solve(mat.NewVecDense(1, []float64{1}), 1e-6, 1, 100)
	suite.Require().NoError(err)

	suite.InDelta(expectedPi, PiStag(suite.node), 1e-5)
}

func (suite CompressorFittingTetSuite) TestPiBiParabolicShape() {
	fg := common.Scaler(ditributions.GetUnitBiParabolic(0, 2, 1, 0.1, 0.1))
	eqSys := GetPiFitEqSys1D(suite.node, expectedPi, suite.setter, fg)

	solver, err := newton.NewUniformNewtonSolver(eqSys, 1e-5, newton.NoLog)
	suite.Require().NoError(err)

	_, err = solver.Solve(mat.NewVecDense(1, []float64{1}), 1e-6, 1, 100)
	suite.Require().NoError(err)

	suite.InDelta(expectedPi, PiStag(suite.node), 1e-5)
}

func (suite CompressorFittingTetSuite) TestPiEtaUnitShape() {
	htFg := common.Scaler(ditributions.GetUnitConstant())
	etaFg := common.Scaler(ditributions.GetUnitConstant())
	eqSys := GetCompressorPiEtaEqSys(suite.node, htFg, expectedPi, etaFg, expectedEta)

	solver, err := newton.NewUniformNewtonSolver(eqSys, 1e-5, newton.NoLog)
	suite.Require().NoError(err)

	_, err = solver.Solve(mat.NewVecDense(2, []float64{0.5, 0.8}), 1e-6, 1, 100)
	suite.Require().NoError(err)

	suite.InDelta(expectedPi, PiStag(suite.node), 1e-5)
	suite.InDelta(expectedEta, EtaStag(suite.node), 1e-5)
}

func TestCompressorFittingTestSuite(t *testing.T) {
	suite.Run(t, new(CompressorFittingTetSuite))
}
