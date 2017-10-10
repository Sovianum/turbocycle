package application

import "strconv"

type SingleCompressorDataPoint struct {
	Pi               float64
	MassRate         float64
	SpecificFuelRate float64
	Efficiency       float64
}

func (point SingleCompressorDataPoint) ToRecord() []string {
	return []string{
		strconv.FormatFloat(point.Pi, 'f', -1, 64),
		strconv.FormatFloat(point.MassRate, 'f', -1, 64),
		strconv.FormatFloat(point.SpecificFuelRate, 'f', -1, 64),
		strconv.FormatFloat(point.Efficiency, 'f', -1, 64),
	}
}
