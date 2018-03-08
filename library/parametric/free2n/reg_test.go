package free2n

import (
	"fmt"
	"testing"

	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	c "github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestNewDoubleShaftRegFreeScheme_Smoke(t *testing.T) {
	var scheme = getUnitRegFreeTestScheme()
	var network, err = scheme.GetNetwork()
	assert.Nil(t, err)

	err = network.Solve(1, 2, 100, 1e-3)
	assert.Nil(t, err)

	var sysCall = variator.SysCallFromNetwork(
		network, scheme.Assembler().GetVectorPort(), 1, 2, 100, 1e-3,
	)
	var variators = scheme.Variators()
	var solverGen = newton.NewUniformNewtonSolverGen(1e-4, newton.DefaultLog)

	var variatorSolver = variator.NewVariatorSolver(sysCall, variators, solverGen)

	_, err = variatorSolver.Solve(
		mat.NewVecDense(6, []float64{1, 1, 0.034, 1, 1, 1}),
		1e-5, 0.5, 10000,
	)
	assert.Nil(t, err)
	fmt.Println(
		scheme.Regenerator().HotInput().TemperatureInput().GetState(),
		scheme.Regenerator().ColdInput().TemperatureInput().GetState(),
		scheme.Regenerator().HotOutput().TemperatureOutput().GetState(),
		scheme.Regenerator().ColdOutput().TemperatureOutput().GetState(),
	)
	fmt.Println(scheme.Compressor().PStagIn(), scheme.Compressor().TStagIn())
	fmt.Println(
		scheme.Regenerator().ColdInput().PressureInput().GetState().Value(),
		scheme.Regenerator().ColdInput().TemperatureInput().GetState().Value(),
	)
	fmt.Println(scheme.Burner().PStagOut(), scheme.Burner().TStagOut())
	fmt.Println(scheme.CompressorTurbine().PStagOut(), scheme.CompressorTurbine().TStagOut())
	fmt.Println(scheme.FreeTurbine().PStagOut(), scheme.FreeTurbine().TStagOut())
}

func getUnitRegFreeTestScheme() DoubleShaftRegFreeScheme {
	return getTestFreeRegScheme(
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(lambdaU, normPiStag float64) float64 {
			return 1
		},
		func(normMassRate, normPiStag float64) float64 {
			return 1
		},
		func(normMassRate, normPiStag float64) float64 {
			return 1
		},
		func(lambda float64) float64 {
			return 1
		},
		func(normRPM float64) float64 {
			return normRPM
		},
	)
}

func getTestFreeRegScheme(
	ctNormMassRateChar, ctNormEtaChar,
	ftNormMassRateChar, ftNormEtaChar c.TurbineCharFunc,
	compressorNormEtaChar, compressorNormRpmChar c.CompressorCharFunc,
	burnerSigmaFunc func(lambda float64) float64,
	powerChar func(normRPM float64) float64,
) DoubleShaftRegFreeScheme {
	var compressor = c.NewParametricCompressorNode(
		cMassRate0, piC0, rpmC0, etaC0, t0, p0, precision,
		compressorNormEtaChar, compressorNormRpmChar,
	)
	var compressorPipe = c.NewPressureLossNode(cpSigma)

	var burner = c.NewParametricBurnerNode(
		fuel.GetCH4(), tFuel, t0, etaBurn, lambdaIn0, p0*piC0, tStagIn0,
		cMassRate0, fuelMassRateRel0, precision, 1, nodes.DefaultN, burnerSigmaFunc,
	)

	var cTurbine = c.NewSimpleParametricTurbineNode(
		ctMassRate0, piCT0, etaCT0, t0, p0, dCTMean, precision,
		0, 0, 0,
		ctNormMassRateChar, ctNormEtaChar,
	)
	var ctPipe = c.NewPressureLossNode(ctpSigma)

	var fTurbine = c.NewSimpleParametricTurbineNode(
		ftMassRate0, piFT0, etaFT0, t0, p0, dFTMean, precision,
		0, 0, 0,
		ftNormMassRateChar, ftNormEtaChar,
	)
	var ftPipe = c.NewPressureLossNode(ftpSigma)

	var payload = c.NewPayload(payloadRPM0, -payloadPower0, powerChar)

	var regenerator = c.NewParametricRegeneratorNode(
		gases.GetAir(), gases.GetAir(), cMassRate0, cMassRate0, tGas*0.75, 500,
		3e5, 3e5, 20, 20, 0.9,
		1e-3, 1.5e-3, 1e-3, 1, nodes.DefaultN,
		c.FrowardTDrop,
		c.DefaultNuFunc, c.DefaultNuFunc,
	)

	var breaker = helper.NewComplexCycleBreakNode(gases.GetAir(), 1300, 1e5, 5.5)

	return NewDoubleShaftRegFreeScheme(
		gases.GetAir(), tAtm, pAtm, tGas,
		0.99, compressor, compressorPipe, regenerator, breaker,
		burner, cTurbine, ctPipe, fTurbine, ftPipe, payload,
	)
}
