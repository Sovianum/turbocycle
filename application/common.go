package application

import (
	"errors"
	"github.com/Sovianum/turbocycle/schemes"
	"strconv"
)

type SingleCompressorScheme interface {
	schemes.Scheme
	schemes.SingleCompressor
}

type SingleCompressorDataPoint struct {
	Pi            float64
	MassRate      float64
	SpecificPower float64
	Efficiency    float64
}

func (point SingleCompressorDataPoint) ToRecord() []string {
	return []string{
		strconv.FormatFloat(point.Pi, 'f', -1, 64),
		strconv.FormatFloat(point.MassRate, 'f', -1, 64),
		strconv.FormatFloat(point.SpecificPower, 'f', -1, 64),
		strconv.FormatFloat(point.Efficiency, 'f', -1, 64),
	}
}

func GetSingleCompressorDataGenerator(
	scheme SingleCompressorScheme, power float64, relaxCoef float64, iterNum int,
) func(pi float64) (SingleCompressorDataPoint, error) {
	return func(pi float64) (SingleCompressorDataPoint, error) {
		scheme.Compressor().SetPiStag(pi)
		var converged, err = scheme.GetNetwork().Solve(relaxCoef, iterNum, 0.001)
		if err != nil {
			return SingleCompressorDataPoint{}, err
		}
		if !converged {
			return SingleCompressorDataPoint{}, errors.New("not converged")
		}

		return SingleCompressorDataPoint{
			Pi:            pi,
			Efficiency:    schemes.GetEfficiency(scheme),
			MassRate:      schemes.GetMassRate(power, scheme),
			SpecificPower: scheme.GetSpecificPower(),
		}, nil
	}
}
