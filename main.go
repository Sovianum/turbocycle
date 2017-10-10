package main

import (
	"os"

	"encoding/csv"

	"github.com/Sovianum/turbocycle/application"
	"github.com/Sovianum/turbocycle/application/two_shafts"
	two_shafts2 "github.com/Sovianum/turbocycle/application/two_shafts_regenerator"
)

const (
	power = 16e6
	relaxCoef = 0.1
	iterNum = 100

	dataRoot = "/home/artem/Documents/University/CoolingSystemProject/notebooks/cycle/data/"
)

func main() {
	if err := saveTwoShaftSchemeData(
		two_shafts.GetInitedTwoShaftsScheme(), 7, 0.1, 120, dataRoot + "2n.csv",
	); err != nil {
		panic(err)
	}

	if err := saveTwoShaftSchemeData(
		two_shafts2.GetInitedTwoShaftsRegeneratorScheme(), 5, 0.1, 100, dataRoot + "2nr.csv",
	); err != nil {
		panic(err)
	}
}

func saveTwoShaftSchemeData(scheme application.SingleCompressorScheme, startPi, piStep float64, stepNum int, filename string) error {
	var piArr []float64

	for i := 0; i != stepNum; i++ {
		piArr = append(piArr, startPi + float64(i) * piStep)
	}

	var records [][]string
	var generator = application.GetSingleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)
	for _, pi := range piArr {
		var point, err = generator(pi)
		if err != nil {
			return err
		}
		records = append(records, point.ToRecord())
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.WriteAll(records)

	return nil
}
