package main

import (
	"os"

	"github.com/Sovianum/turbocycle/application/two_shafts"
	"encoding/csv"
)

func main() {
	var scheme = two_shafts.GetInitedTwoShaftsScheme()
	var piArr []float64

	var startPi float64 = 7
	var piStep float64 = 0.2
	for i := 0; i != 200; i++ {
		piArr = append(piArr, startPi + float64(i) * piStep)
	}

	var records [][]string
	var generator = two_shafts.GetDataGenerator(scheme, 16e6, 0.1, 100)
	for _, pi := range piArr {
		var point, err = generator(pi)
		if err != nil {
			panic(err)
		}
		records = append(records, point.ToRecord())
	}

	file, err := os.Create("/tmp/result.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.WriteAll(records)
}
