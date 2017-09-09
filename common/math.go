package common

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

func Converged(x0, x1, precision float64) bool {
	return math.Abs(x0-x1)/math.Abs(x0) <= precision
}

func GetRelResidual(x0 float64, x1 float64) float64 {
	return math.Abs(x0-x1) / math.Abs(x0)
}

func Lerp(x0 float64, x1 float64, interpCoef float64) float64 {
	return x0*(1-interpCoef) + x1*interpCoef
}

func Interp(x float64, xArr []float64, yArr []float64) (float64, error) {
	if len(xArr) != len(yArr) {
		return 0, errors.New(fmt.Sprintf("len(xArr) = %d, len(yArr) = %d", len(xArr), len(yArr)))
	}

	if !sort.SliceIsSorted(xArr, func(i, j int) bool {
		return xArr[i] < xArr[j]
	}) {
		return 0, errors.New("xArr no sorted")
	}

	if x < xArr[0] {
		return 0, errors.New(fmt.Sprintf("x(x == %f) < xArr[0](xArr[0] == %f)", x, xArr[0]))
	}
	if x > xArr[len(xArr)-1] {
		return 0, errors.New(fmt.Sprintf("x(x == %f) > xArr[-1](xArr[-1] == %f)", x, xArr[len(xArr)-1]))
	}

	for i := 0; i != len(xArr)-1; i++ {
		if x <= xArr[i] && x < xArr[i+1] {
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

	return result * dx / (x1 - x0)
}

func ApproxEqual(expected, got, precision float64) bool {
	var min = expected * (1 - precision/2)
	var max = expected * (1 + precision/2)

	return min <= got && got <= max
}
