package main

import (
	"os"

	"encoding/csv"

	"github.com/Sovianum/turbocycle/application"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/application/three_shafts_cool_regenerator"
)

const (
	power = 16e6
	relaxCoef = 0.1
	iterNum = 100

	dataRoot = "/home/artem/Documents/University/CoolingSystemProject/notebooks/cycle/data/"
)

func main() {
	//if err := saveTwoShaftSchemeData(
	//	two_shafts.GetInitedTwoShaftsScheme(), 7, 0.1, 120, dataRoot + "2n.csv",
	//); err != nil {
	//	panic(err)
	//}
	//
	//if err := saveTwoShaftSchemeData(
	//	two_shafts_regenerator.GetInitedTwoShaftsRegeneratorScheme(), 5, 0.1, 100, dataRoot + "2nr.csv",
	//); err != nil {
	//	panic(err)
	//}
	//
	//if err := saveThreeShaftsSchemeData(
	//	three_shafts.GetInitedThreeShaftsScheme(),
	//	8, 0.1, 120,
	//	0.15, 0.1, 8,
	//	dataRoot + "3n.csv",
	//); err != nil {
	//	panic(err)
	//}
	//
	//if err := saveThreeShaftsSchemeData(
	//	three_shafts_regenerator.GetInitedThreeShaftsRegeneratorScheme(),
	//	7, 0.1, 120,
	//	0.15, 0.1, 8,
	//	dataRoot + "3nr.csv",
	//); err != nil {
	//	panic(err)
	//}

	//if err := saveThreeShaftsSchemeData(
	//	three_shafts_cool.GetInitedThreeShaftsCoolingScheme(),
	//	20, 0.1, 120,
	//	0.15, 0.1, 8,
	//	dataRoot + "3nc.csv",
	//); err != nil {
	//	panic(err)
	//}

	if err := saveThreeShaftsSchemeData(
		three_shafts_cool_regenerator.GetInitedThreeShaftsCoolRegeneratorScheme(),
		8, 0.1, 150,
		0.15, 0.1, 8,
		dataRoot + "3ncr.csv",
	); err != nil {
		panic(err)
	}

}

func saveThreeShaftsSchemeData(
	scheme schemes.ThreeShaftsScheme,
	startPi, piStep float64, piStepNum int,
	startPiFactor, piFactorStep float64, piFactorStepNum int,
	filename string,
) error {
	var piArr []float64
	for i := 0; i != piStepNum; i++ {
		piArr = append(piArr, startPi + float64(i) * piStep)
	}

	var piFactorArr []float64
	for i := 0; i != piFactorStepNum; i++ {
		piFactorArr = append(piFactorArr, startPiFactor+ float64(i) *piFactorStep)
	}

	var records [][]string
	var generator = application.GetDoubleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)
	for _, piFactor := range piFactorArr {
		for _, pi := range piArr {
			var point, err = generator(pi, piFactor)
			if err != nil {
				return err
			}
			records = append(records, point.ToRecord())
		}
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
