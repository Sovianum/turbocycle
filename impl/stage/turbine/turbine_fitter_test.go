package turbine

import (
	"fmt"
	"testing"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"github.com/Sovianum/turbocycle/impl/stage/ditributions"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

const (
	expectedPi  = 3
	expectedEta = 0.87
)

type TurbineFittingTetSuite struct {
	suite.Suite
	node       StagedTurbineNode
	htVariator variator.Variator
	piSolver   variator.VariatorSolver
}

func (suite *TurbineFittingTetSuite) SetupTest() {
	suite.node = getTestTurbine()
	suite.Require().NoError(suite.node.Process())
	suite.htVariator = variator.FromCallables(suite.node.Ht, suite.node.SetHt)
	suite.piSolver = variator.NewVariatorSolver(
		func() (*mat.VecDense, error) {
			if err := suite.node.Process(); err != nil {
				return nil, err
			}
			return mat.NewVecDense(1, []float64{PiStag(suite.node) - expectedPi}), nil
		},
		[]variator.Variator{suite.htVariator},
		newton.NewUniformNewtonSolverGen(1e-5, newton.NoLog),
	)
	suite.Require().NoError(suite.node.Process())
}

func (suite TurbineFittingTetSuite) TestFitPi() {
	_, err := suite.piSolver.Solve(mat.NewVecDense(1, []float64{suite.node.Ht()}), 1e-6, 0.1, 1000)
	suite.Require().NoError(err)

	suite.InDelta(expectedPi, PiStag(suite.node), 1e-5)
}

func (suite TurbineFittingTetSuite) TestFitPiAndEta() {
	fmt.Println(EtaStag(suite.node))
	eqSys := GetTurbinePiEtaEqSys(
		suite.node,
		common.Scaler(ditributions.GetUnitConstant().Scale(phi)),
		common.Scaler(ditributions.GetUnitConstant().Scale(psi)),
		expectedPi, expectedEta,
	)
	solver, err := newton.NewUniformNewtonSolver(eqSys, 1e-5, newton.NoLog)
	suite.Require().NoError(err)

	_, err = solver.Solve(mat.NewVecDense(2, []float64{suite.node.Ht(), phi}), 1e-6, 1, 100)
	suite.Require().NoError(err)

	suite.InDelta(expectedPi, PiStag(suite.node), 1e-5)
	suite.InDelta(expectedEta, EtaStag(suite.node), 1e-5)
}

func TestCompressorFittingTestSuite(t *testing.T) {
	suite.Run(t, new(TurbineFittingTetSuite))
}
