package gases

import (
	"github.com/Sovianum/turbocycle/common"
)

func GetAir() Gas {
	return air{}
}

type air struct {}

func (air) Cp(t float64) float64 {
	var tArr = []float64{
		260, 333,  393,  413,  433,  453,  473,  523,  573,  623,  673,  773, 873,  973, 1073, 1173, 1273, 1373, 1473,
	}
	var cpArr = []float64{
		1000, 1005, 1009, 1013, 1017, 1022, 1026, 1038, 1047, 1059, 1068, 1093, 1114, 1135, 1156, 1172, 1185, 1197, 1210,
	}

	var cp, err = common.Interp(t, tArr, cpArr)
	if err != nil {
		panic(err)
	}

	return cp
}

func (air) R() float64 {
	return 287.
}

