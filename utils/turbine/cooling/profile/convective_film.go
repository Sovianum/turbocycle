package profile

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode/boundary"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode/forward"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"gonum.org/v1/gonum/mat"
)

const (
	defaultN = 100
)

func NewSlitInfo(coord, thickness, area, velocityCoef, massRateCoef float64) SlitInfo {
	return SlitInfo{
		Coord:        coord,
		Thickness:    thickness,
		Area:         area,
		VelocityCoef: velocityCoef,
		MassRateCoef: massRateCoef,

		thermoPoint: thermoPoint{},
		activated:   false,
	}
}

type SlitInfo struct {
	Coord        float64
	Thickness    float64
	Area         float64
	VelocityCoef float64
	MassRateCoef float64

	thermoPoint thermoPoint
	activated   bool
}

type thermoPoint struct {
	coord    float64
	theta    float64
	gasTemp  float64
	massRate float64
}

func NewConvFilmTemperatureSystem(
	solver forward.Solver,
	coolerMassRate0 float64,
	cooler gases.Gas,
	gas gases.Gas,
	gasTStag, gasPStag, coolerPStag func(x float64) float64,
	lambdaFunc func(xRel float64) float64,
	alphaCoolerConv, alphaGasConv cooling.AlphaLaw,
	slitInfoArray []SlitInfo,
	wallThk, lambdaM func(t float64) float64,
	segment geom.Segment,
) TemperatureSystem {
	return &convFilmTemperatureSystem{
		solver:          solver,
		coolerMassRate0: coolerMassRate0,

		cooler: cooler,
		gas:    gas,

		gasTempStag:        gasTStag,
		gasPressureStag:    gasPStag,
		coolerPressureStag: coolerPStag,

		lambdaFunc: lambdaFunc,

		alphaCoolerConv: alphaCoolerConv,
		alphaGasConv:    alphaGasConv,

		slitInfoArray: slitInfoArray,

		wallThk: wallThk,
		lambdaM: lambdaM,

		segment:       segment,
		segmentLength: geom.ApproxLength(segment, 0, 1, defaultN),

		lengthParameter: 0,
		lengthMassRate:  coolerMassRate0,
	}

}

type convFilmTemperatureSystem struct {
	solver         forward.Solver
	wallTempSolver boundary.Solver

	coolerMassRate0 float64

	cooler gases.Gas
	gas    gases.Gas

	gasTempStag        func(x float64) float64
	gasPressureStag    func(x float64) float64
	coolerPressureStag func(x float64) float64

	lambdaFunc func(xRel float64) float64

	alphaCoolerConv cooling.AlphaLaw
	alphaGasConv    cooling.AlphaLaw

	slitInfoArray []SlitInfo

	wallThk func(x float64) float64
	lambdaM func(t float64) float64

	segment       geom.Segment
	segmentLength float64

	solutionStep float64

	lengthParameter float64
	lengthMassRate  float64
}

func (system *convFilmTemperatureSystem) Solve(t0, theta0, tMax, maxStep float64) TemperatureSolution {
	system.lengthParameter = 0
	//system.lengthMassRate = system.coolerMassRate0
	system.lengthMassRate = system.coolerMassRate0 / 2

	var solution = system.solver.Solution(system.dThetaDX, t0, theta0, tMax, maxStep)
	system.solutionStep = solution.Step()

	var tArr, tAirArr = solution.Build()

	var solutionLen = len(tArr)
	var tSolution = NewTemperatureSolution(solutionLen)
	tSolution.ParametricCoord = tArr
	tSolution.AirTemperature = tAirArr

	system.extendSolutionArray(&tSolution, t0, system.solutionStep, solutionLen)

	system.initWalTempSolver(tSolution.LengthCoord, tSolution.AirTemperature, tSolution.FilmTemperature)
	if wallTempSolution, err := system.wallTempSolver.Solve(); err != nil {
		panic(err)
	} else {
		_, t := wallTempSolution.Build()
		tSolution.SmoothWallTemperature = t
	}

	return tSolution
}

func (system *convFilmTemperatureSystem) CpAir(theta float64) float64 {
	return system.cooler.Cp(theta)
}

func (system *convFilmTemperatureSystem) GasTemp(lengthCoord float64) float64 {
	return system.gasTempStag(lengthCoord)
}

func (system *convFilmTemperatureSystem) AlphaAir(lengthCoord, theta float64) float64 {
	return system.alphaCoolerConv(lengthCoord, theta)
}

func (system *convFilmTemperatureSystem) AlphaGas(lengthCoord, theta float64) float64 {
	return system.alphaFilm(lengthCoord, theta)
}

func (system *convFilmTemperatureSystem) K(x, theta float64) float64 {
	return system.heatTransferCoef(x, theta)
}

func (system *convFilmTemperatureSystem) extendSolutionArray(
	solution *TemperatureSolution,
	tStart, step float64, pointNum int,
) {
	var currLengthCoord float64 = 0
	var currT float64 = 0

	for i := 0; i != pointNum; i++ {
		var point = system.segment.Point(currT)
		solution.X[i] = point.At(0, 0)
		solution.Y[i] = point.At(1, 0)
		solution.LengthCoord[i] = currLengthCoord
		solution.AlphaAir[i] = system.AlphaAir(currLengthCoord, solution.AirTemperature[i])
		solution.AlphaGas[i] = system.AlphaGas(currLengthCoord, solution.AirTemperature[i])
		solution.WallTemperature[i] = system.wallTemp(currLengthCoord, solution.AirTemperature[i])
		solution.FilmTemperature[i] = system.multiSlitFilmTemperature(currLengthCoord)
		solution.HeatTransferCoef[i] = system.heatTransferCoef(currLengthCoord, solution.AirTemperature[i])
		solution.FilmEfficiency[i] =
			(system.gasTempStag(currLengthCoord) - solution.FilmTemperature[i]) /
				(system.gasTempStag(currLengthCoord) - solution.AirTemperature[i])

		var segmentDerivative = geom.ApproxDerivative1(system.segment, currT, step)
		var lengthStep = mat.Norm(segmentDerivative, 2) * step

		currT += step
		currLengthCoord += lengthStep
	}

	for i, info := range system.slitInfoArray {
		solution.SlitsSolution = append(
			solution.SlitsSolution,
			SlitSolution{
				Id:                i + 1,
				SlitInfo:          info,
				Velocity:          system.coolerSlitVelocity(info.Coord, info.thermoPoint.theta, info.VelocityCoef),
				BlowingParameter:  system.slitBlowingParameter(info.Coord, info.thermoPoint.theta, info.VelocityCoef),
				TemperatureFactor: info.thermoPoint.theta / system.gasTempStag(info.Coord),
				MassRate:          info.thermoPoint.massRate,
				MassRateRel:       info.thermoPoint.massRate / system.coolerMassRate0,
			},
		)
	}
}

func (system *convFilmTemperatureSystem) dThetaDX(t, theta float64) float64 {
	// slits are activated when we get cooler temperature at this particular point
	system.activateSlits(system.lengthParameter, theta, system.gasTempStag(system.lengthParameter))

	var segmentDerivative = geom.ApproxDerivative1(system.segment, t, system.solutionStep)
	var geomFactor = mat.Norm(segmentDerivative, 2)

	var coolerCp = system.cooler.Cp(theta)
	var massRateFactor = 1 / (system.lengthMassRate * coolerCp)

	var k = system.heatTransferCoef(system.lengthParameter, theta)
	var filmTemperature = system.multiSlitFilmTemperature(system.lengthParameter)

	var result = geomFactor * massRateFactor * k * (filmTemperature - theta)

	// side effect is used to prevent length calculation on each step
	system.lengthParameter += geomFactor * system.solutionStep
	return result
}

func (system *convFilmTemperatureSystem) activateSlits(lengthCoord, theta, gasTemp float64) {
	for i := 0; i != len(system.slitInfoArray); i++ {
		if !system.slitInfoArray[i].activated && system.slitInfoArray[i].Coord <= lengthCoord {
			system.slitInfoArray[i].activated = true

			system.slitInfoArray[i].thermoPoint = thermoPoint{
				theta:   theta,
				coord:   system.slitInfoArray[i].Coord,
				gasTemp: gasTemp,
			}
			// slit mass rate is calculated after initialization in order to pass correct
			// parameters to mass rate calculating function
			var slitMassRate = system.coolerMassRateSlit(system.slitInfoArray[i])
			system.slitInfoArray[i].thermoPoint.massRate = slitMassRate
			// decrease cooler mass rate by value which flows to the slit
			system.lengthMassRate -= slitMassRate
		}
	}
}

func (system *convFilmTemperatureSystem) wallTemp(lengthCoord, theta float64) float64 {
	var tFilm = system.multiSlitFilmTemperature(lengthCoord)
	var alphaFilm = system.alphaFilm(lengthCoord, theta)

	var kFactor = system.heatTransferCoef(lengthCoord, theta) / alphaFilm
	var tFactor = tFilm - theta
	return tFilm - kFactor*tFactor
}

func (system *convFilmTemperatureSystem) initWalTempSolver(lengthCoord, tAir, tFilm []float64) {
	fArr := make([]float64, len(lengthCoord))
	gArr := make([]float64, len(lengthCoord))
	hArr := make([]float64, len(lengthCoord))

	for i, xi := range lengthCoord {
		wallThkVal := system.wallThk(xi)
		lambdaVal := system.lambdaM(xi)
		tFilmVal := tFilm[i]
		tAirVal := tAir[i]
		alphaAirVal := system.AlphaAir(xi, tAirVal)
		alphaFilmVal := system.AlphaGas(xi, tFilmVal)

		//fArr[i] == 0 for all i
		gArr[i] = -(alphaAirVal + alphaFilmVal) / (wallThkVal * lambdaVal)
		hArr[i] = -(alphaAirVal*tAirVal + alphaFilmVal*tFilmVal) / (wallThkVal * lambdaVal)
	}

	pointNum := len(lengthCoord)
	lastStep := lengthCoord[pointNum-1] - lengthCoord[pointNum-2]
	dtAir := tAir[pointNum-1] - tAir[pointNum-2]
	dtFilm := tFilm[pointNum-1] - tFilm[pointNum-2]

	system.wallTempSolver = boundary.NewSolverFromArrays(
		lengthCoord, fArr, gArr, hArr,
		boundary.NewSecondTypeBC(0),
		boundary.NewSecondTypeBC((dtAir+dtFilm*0)/(2*lastStep)),
		//boundary.NewSecondTypeBC(0),
	)
}

func (system *convFilmTemperatureSystem) heatTransferCoef(lengthCoord, theta float64) float64 {
	var alphaCooler = system.alphaCoolerConv(lengthCoord, theta)
	var alphaFilm = system.alphaFilm(lengthCoord, theta)
	var delta = system.wallThk(lengthCoord)
	var lambdaM = system.lambdaM(theta)

	var enom = alphaCooler
	var denom = 1 + alphaCooler*(1/alphaFilm+delta/lambdaM)
	return enom / denom
}

func (system *convFilmTemperatureSystem) alphaFilm(x, theta float64) float64 {
	var alphaConv = system.alphaGasConv(x, theta)
	var filmFactor = system.alphaFilmFactor(x, system.slitInfoArray)
	return alphaConv * filmFactor
}

func (system *convFilmTemperatureSystem) alphaFilmFactor(x float64, slitInfoArray []SlitInfo) float64 {
	var result float64 = 1
	for _, info := range slitInfoArray {
		if info.Coord <= x && info.activated {
			var blowParameter = system.slitBlowingParameter(info.Coord, info.thermoPoint.theta, info.VelocityCoef)
			var xRel = (x - info.Coord) / info.Thickness
			result *= 1 + 2*blowParameter/xRel
		}
	}
	if result > 100 {
		return 100
	}
	return result
}

func (system *convFilmTemperatureSystem) multiSlitFilmTemperature(x float64) float64 {
	var filmEfficiencyFactors = make([]float64, len(system.slitInfoArray))
	var complementaryFactors = make([]float64, len(system.slitInfoArray))

	for i, info := range system.slitInfoArray {
		var parameter = system.filmEfficiencyParameter(x, info)
		var factor = system.filmEfficiencyFactor(parameter)

		filmEfficiencyFactors[i] = factor
		complementaryFactors[i] = 1 - factor
	}

	var term1 = system.gasTempStag(x) * common.Product(complementaryFactors)

	var term2 float64 = 0
	for i := range filmEfficiencyFactors {
		if system.slitInfoArray[i].activated {
			term2 += filmEfficiencyFactors[i] *
				system.slitInfoArray[i].thermoPoint.theta *
				common.Product(complementaryFactors[i+1:])
		}
	}

	return term1 + term2
}

func (system *convFilmTemperatureSystem) filmEfficiencyFactor(filmEfficiencyParameter float64) float64 {
	if filmEfficiencyParameter < 0 {
		return 0
	} else if 0 <= filmEfficiencyParameter && filmEfficiencyParameter < 3 {
		return 1
	} else if 3 <= filmEfficiencyParameter && filmEfficiencyParameter < 11 {
		return math.Pow(filmEfficiencyParameter/3, -0.285)
	} else {
		return math.Pow(filmEfficiencyParameter/7.43, -0.95)
	}
}

func (system *convFilmTemperatureSystem) filmEfficiencyParameter(x float64, slitInfo SlitInfo) float64 {
	if !slitInfo.activated {
		return -1 // means information about thermal parameters in slitPoint is unknown
	}

	var reFactor = math.Pow(system.gasSlitRe(slitInfo.Coord, slitInfo.Thickness), -0.25)
	var blowingFactor = math.Pow(system.slitBlowingParameter(
		slitInfo.Coord,
		slitInfo.thermoPoint.theta,
		slitInfo.VelocityCoef,
	), -1.3)
	var tempFactor = math.Pow(slitInfo.thermoPoint.theta/system.gasTempStag(x), -1.25)
	var coordFactor = (x - slitInfo.Coord) / slitInfo.Thickness

	return reFactor * blowingFactor * tempFactor * coordFactor
}

func (system *convFilmTemperatureSystem) gasSlitRe(x, slitThk float64) float64 {
	var tStat = system.gasTempStat(x)

	var density = gases.Density(system.gas, tStat, system.gasPressureStag(x))
	var velocity = system.gasVelocity(x)
	var viscosity = system.gas.Mu(tStat)

	return velocity * density * slitThk / viscosity
}

func (system *convFilmTemperatureSystem) coolerMassRate(x float64) float64 {
	var result = system.coolerMassRate0
	for _, info := range system.slitInfoArray {
		if info.activated && info.Coord < x {
			result -= info.thermoPoint.massRate
		}
	}
	return result
}

func (system *convFilmTemperatureSystem) coolerMassRateSlit(slitInfo SlitInfo) float64 {
	var gasPressure = system.gasPressureStat(slitInfo.Coord)

	var coolerK = gases.K(system.cooler, slitInfo.thermoPoint.theta)
	var coolerDensity = system.coolerSlitDensityStat(
		slitInfo.Coord, slitInfo.thermoPoint.theta, slitInfo.VelocityCoef,
	)
	var coolerPressure = system.coolerPressureStag(slitInfo.Coord)

	var pi = gasPressure / coolerPressure

	var kFactor = 2 * coolerK / (coolerK - 1)
	var dimFactor = coolerPressure * coolerDensity
	var piFactor = math.Pow(pi, 2/coolerK) * (1 - math.Pow(pi, (coolerK-1)/coolerK))

	var result = slitInfo.Area * slitInfo.MassRateCoef * math.Sqrt(kFactor*dimFactor*piFactor)

	if result > system.lengthMassRate {
		panic("not enough mass rate")
	}
	return result
}

func (system *convFilmTemperatureSystem) slitBlowingParameter(slitCoord, slitTheta, velocityCoef float64) float64 {
	var gasVelocity = system.gasVelocity(slitCoord)
	var gasDensity = system.gasDensityStat(slitCoord)

	var coolerVelocity = system.coolerSlitVelocity(slitCoord, slitTheta, velocityCoef)
	var coolerDensity = system.coolerSlitDensityStat(slitCoord, slitTheta, velocityCoef)

	var result = coolerDensity * coolerVelocity / (gasDensity * gasVelocity)
	return result
}

func (system *convFilmTemperatureSystem) coolerSlitDensityStat(slitCoord, slitTheta, velocityCoef float64) float64 {
	var lambda = system.airSlitLambda(slitCoord, slitTheta, velocityCoef)
	var densityStag = system.coolerDensityStag(slitCoord, slitTheta)
	var k = gases.K(system.cooler, slitTheta)

	return densityStag * gdf.Epsilon(lambda, k)
}

func (system *convFilmTemperatureSystem) airSlitLambda(slitCoord, slitTheta, velocityCoef float64) float64 {
	var coolerSlitVelocity = system.coolerSlitVelocity(slitCoord, slitTheta, velocityCoef)
	var coolerCritSpeedSound = system.coolerCritSpeedSound(slitCoord, slitTheta)
	return coolerSlitVelocity / coolerCritSpeedSound
}

func (system *convFilmTemperatureSystem) coolerSlitVelocity(slitCoord, slitTheta, velocityCoef float64) float64 {
	var kAir = gases.K(system.cooler, slitTheta)
	var kFactor = 2 * kAir / (kAir - 1)

	var rAir = system.cooler.R()
	var tempFactor = rAir * slitTheta

	var pAir = system.coolerPressureStag(slitCoord)
	var pGas = system.gasPressureStat(slitCoord)
	var pFactor = 1 - math.Pow(pGas/pAir, (kAir-1)/kAir)

	var velocity = velocityCoef * math.Sqrt(kFactor*tempFactor*pFactor)
	return velocity
}

func (system *convFilmTemperatureSystem) coolerCritSpeedSound(x, theta float64) float64 {
	var k = gases.K(system.cooler, theta)
	var r = system.cooler.R()
	return gdf.ACrit(k, r, theta)
}

func (system *convFilmTemperatureSystem) coolerDensityStag(x float64, theta float64) float64 {
	return gases.Density(system.cooler, theta, system.gasPressureStag(x))
}

func (system *convFilmTemperatureSystem) gasDensityStat(x float64) float64 {
	var tStat = system.gasTempStat(x)
	var pStat = system.gasPressureStat(x)
	var densityStat = gases.Density(system.gas, tStat, pStat)
	return densityStat
}

func (system *convFilmTemperatureSystem) gasPressureStat(x float64) float64 {
	var tStag = system.gasTempStag(x)
	var tStat = system.gasTempStat(x)

	var tRelation = tStag / tStat
	var kGas = gases.KMean(system.gas, tStag, tStag, defaultN)

	var tDivisor = math.Pow(tRelation, kGas/(kGas-1))
	return system.gasPressureStag(x) / tDivisor
}

func (system *convFilmTemperatureSystem) gasTempStat(x float64) float64 {
	var tGas = system.gasTempStag(x)
	var velocity = system.gasVelocity(x)
	var cp = system.gas.Cp(tGas)
	var tStat = tGas - velocity*velocity/(2*cp)
	return tStat
}

func (system *convFilmTemperatureSystem) gasVelocity(x float64) float64 {
	var lambda = system.lambda(x)
	var critSpeedSound = system.gasCritSpeedSound(x)
	var velocity = lambda * critSpeedSound
	return velocity
}

func (system *convFilmTemperatureSystem) gasCritSpeedSound(x float64) float64 {
	var tStag = system.gasTempStag(x)
	return gdf.ACrit(
		gases.K(system.gas, tStag),
		system.gas.R(),
		tStag,
	)
}

func (system *convFilmTemperatureSystem) lambda(x float64) float64 {
	return system.lambdaFunc(x / system.segmentLength)
}
