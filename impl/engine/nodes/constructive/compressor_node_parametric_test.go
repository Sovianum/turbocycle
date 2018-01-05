package constructive

import (
	"fmt"
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

const (
	massRate0 = 40
	piC0      = 10
	rpm0      = 10000
	etaAd0    = 0.86
)

func TestParametricCompressorNode_Process_Smoke_Constant(t *testing.T) {
	var compressor = getTestParametricCompressor(
		func(float64, float64) float64 { return 1. },
		func(float64, float64) float64 { return 1. },
	)

	var pi float64 = piC0 / 2
	compressor.SetPiStag(pi)

	var err = compressor.Process()
	assert.Nil(t, err)

	assert.Equal(t, pA*pi, compressor.PStagOut())

	var expectedTemperature = tA * (1 + 1/etaAd0*(math.Pow(pi, (kAir-1)/kAir)-1))
	assert.True(
		t,
		common.ApproxEqual(expectedTemperature, compressor.TStagOut(), 0.01),
		fmt.Sprintf("Expected %f, got %f", expectedTemperature, compressor.TStagOut()),
	)

	var cp = gases.CpMean(gases.GetAir(), compressor.TStagIn(), compressor.TStagOut(), nodes.DefaultN)
	var expectedLabour = cp * (compressor.TStagOut() - compressor.TStagIn())
	assert.True(
		t,
		common.ApproxEqual(expectedLabour, compressor.LSpecific(), 0.01),
		fmt.Sprintf("Expected %f, got %f", expectedLabour, compressor.LSpecific()),
	)

	assert.InDelta(t, rpm0, compressor.RPM(), 1e-7)
}

func TestParametricCompressorNode_Process_Smoke_NonLinear(t *testing.T) {
	var f = func(x float64) float64 { return 1 - 0.1*(1-x)*(1-x) }

	var compressor = getTestParametricCompressor(
		func(normMassRate float64, normPiStag float64) float64 {
			return f(normMassRate)
		},
		func(normMassRate float64, normPiStag float64) float64 { return 1. },
	)

	var massRateVariator = variator.FromCallables(
		compressor.NormMassRate, compressor.SetNormMassRate,
	)

	var massRate = 45.
	var systemCall = func() (*mat.VecDense, error) {
		compressor.Process()
		return mat.NewVecDense(1, []float64{compressor.MassRate() - massRate}), nil
	}

	var solverGen = newton.NewUniformNewtonSolverGen(1e-4)
	var variatorSolver = variator.NewVariatorSolver(systemCall, []variator.Variator{massRateVariator}, solverGen)

	var _, err = variatorSolver.Solve(mat.NewVecDense(1, []float64{1}), 1e-8, 100)
	assert.Nil(t, err)

	assert.InDelta(t, compressor.MassRate(), massRate, 1e-7)
}

func getTestParametricCompressor(normEtaChar, normRpmChar CompressorCharFunc) ParametricCompressorNode {
	var compressor = NewParametricCompressorNode(
		massRate0, piC0, rpm0, etaAd0,
		states.StandardTemperature, states.StandardPressure, 0.01,
		normEtaChar, normRpmChar,
	)
	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gases.GetAir()), states.NewTemperaturePortState(tA), states.NewPressurePortState(pA),
		},
		[]graph.Port{
			compressor.GasInput(), compressor.TemperatureInput(), compressor.PressureInput(),
		},
	)
	return compressor
}
