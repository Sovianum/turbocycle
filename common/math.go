package common

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

func ReduceFloat64(xArr []float64, reducer func(x, y float64) float64, init float64) float64 {
	result := init
	for _, x := range xArr {
		result = reducer(x, result)
	}
	return result
}

func MaxID(xArr []float64) int {
	result := -1
	max := -1e10
	for i, x := range xArr {
		if x > max {
			max = x
			result = i
		}
	}
	return result
}

func Sum(xArr []float64) float64 {
	return ReduceFloat64(xArr, func(x, y float64) float64 { return x + y }, 0)
}

func Product(xArr []float64) float64 {
	return ReduceFloat64(xArr, func(x, y float64) float64 { return x * y }, 1)
}

func SolveIteratively(
	f func(xCurr float64) (xNew float64, err error),
	x0, precision, relaxCoef float64, iterLimit int,
) (float64, error) {
	return SolveIterativelyWithValidation(f, func(float64) error {
		return nil
	}, x0, precision, relaxCoef, iterLimit)
}

func SolveIterativelyWithValidation(
	f func(xCurr float64) (xNew float64, err error),
	validator func(x float64) error,
	x0, precision, relaxCoef float64, iterLimit int,
) (float64, error) {
	fixedPrecision := precision * relaxCoef // precision is increased to mitigate relax coef effect
	interp := func(x1, x2 float64) float64 {
		return relaxCoef*x2 + (1-relaxCoef)*x1
	}

	xCurr := x0
	xNew, err := f(xCurr)
	if err != nil {
		return 0, err
	} else {
		xNew = interp(xCurr, xNew)
	}

	i := 0
	for i < iterLimit && !Converged(xCurr, xNew, fixedPrecision) {
		i++

		xCurr = xNew
		xNew, err = f(xCurr)
		if err != nil {
			return 0, err
		} else {
			xNew = interp(xCurr, xNew)
		}
		if err := validator(xNew); err != nil {
			return 0, fmt.Errorf("validation failed: %s", err.Error())
		}
	}

	if Converged(xCurr, xNew, fixedPrecision) {
		return xNew, nil
	}
	return 0, fmt.Errorf("failed to converge (SolveIteratively)")
}

func Converged(x0, x1, precision float64) bool {
	return math.Abs(x0-x1)/math.Abs(x0) <= precision
}

func GetRelResidual(x0 float64, x1 float64) float64 {
	return math.Abs(x0-x1) / math.Abs(x0)
}

func Lerp(x0 float64, x1 float64, interpCoef float64) float64 {
	return x0*(1-interpCoef) + x1*interpCoef
}

func InterpTolerate(x float64, xArr []float64, yArr []float64) float64 {
	var result, _ = Interp(x, xArr, yArr, false)
	return result
}

func Interp(x float64, xArr []float64, yArr []float64, breakOutOfRange bool) (float64, error) {
	if len(xArr) != len(yArr) {
		return 0, fmt.Errorf("len(xArr) = %d, len(yArr) = %d", len(xArr), len(yArr))
	}

	if !sort.SliceIsSorted(xArr, func(i, j int) bool {
		return xArr[i] < xArr[j]
	}) {
		return 0, errors.New("xArr no sorted")
	}

	if x < xArr[0] {
		if breakOutOfRange {
			return 0, fmt.Errorf("x(x == %f) < xArr[0](xArr[0] == %f)", x, xArr[0])
		}
		x = xArr[0]
	}
	if x > xArr[len(xArr)-1] {
		if breakOutOfRange {
			return 0, fmt.Errorf("x(x == %f) > xArr[-1](xArr[-1] == %f)", x, xArr[len(xArr)-1])
		}
		x = xArr[len(xArr)-1]
	}

	for i := 0; i != len(xArr)-1; i++ {
		if x >= xArr[i] && x < xArr[i+1] {
			var interpCoef = (x - xArr[i]) / (xArr[i+1] - xArr[i])
			return Lerp(yArr[i], yArr[i+1], interpCoef), nil
		}
	}

	return yArr[len(yArr)-1], nil
}

func Average(f func(float64) float64, x0 float64, x1 float64, n int) float64 {
	if x0 == x1 {
		return f(x0)
	}

	return Integrate(f, x0, x1, n) / (x1 - x0)
}

func Integrate(f func(float64) float64, x0 float64, x1 float64, n int) float64 {
	if x0 == x1 {
		return 0
	}

	var dx = (x1 - x0) / float64(n)

	var xVals = make([]float64, n+1)
	for i := 0; i != n+1; i++ {
		xVals[i] = x0 + float64(i)*dx
	}

	var fVals = make([]float64, n)
	for i := 0; i != n; i++ {
		fVals[i] = f((xVals[i] + xVals[i+1]) / 2)
	}

	var result float64 = 0
	for _, val := range fVals {
		result += val
	}

	return result * dx
}

func ApproxEqual(expected, got, precision float64) bool {
	var min = expected - math.Abs(expected)*precision
	var max = expected + math.Abs(expected)*precision

	return min <= got && got <= max
}

func ToDegrees(x float64) float64 {
	return x * 180 / math.Pi
}

func ToRadians(x float64) float64 {
	return x * math.Pi / 180
}

func RoundInt(x float64) int {
	var deltaFloor = math.Abs(x - float64(int(x)))
	var deltaCeil = 1 - deltaFloor

	if x >= 0 {
		if deltaCeil >= deltaFloor {
			return int(x)
		}
		return int(x) + 1
	} else {
		if deltaCeil >= deltaFloor {
			return int(x)
		}
		return int(x) - 1
	}
}

func Factorial(x int) int {
	var result = 1
	for i := 2; i <= x; i++ {
		result *= i
	}
	return result
}

func Arange(start, step float64, stepNum int) []float64 {
	result := make([]float64, stepNum)
	for i := range result {
		result[i] = start + step*float64(i)
	}
	return result
}

func LinSpace(x1, x2 float64, n int) []float64 {
	var step = (x2 - x1) / float64(n-1)
	var result = make([]float64, n)

	for i := 0; i != n; i++ {
		result[i] = x1 + step*float64(i)
	}

	return result
}

func LinScale(x, xFromStart, xFromEnd, xToStart, xToEnd float64) float64 {
	return (xToEnd-xToStart)/(xFromEnd-xFromStart)*(x-xFromStart) + xToStart
}

func Cross2(v1 mat.Vector, v2 mat.Vector) float64 {
	if v1.Len() != 2 {
		panic(fmt.Errorf("first argument has len = %d != 2", v1.Len()))
	}
	if v2.Len() != 2 {
		panic(fmt.Errorf("second argument has len = %d != 2", v1.Len()))
	}

	return v1.At(0, 0)*v2.At(1, 0) - v2.At(0, 0)*v1.At(1, 0)
}
