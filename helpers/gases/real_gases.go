package gases

import (
	"github.com/Sovianum/turbocycle/common"
)

func GetAir() Gas {
	return air{}
}

func GetNitrogen() Gas {
	return nitrogen{}
}

func GetCO2() Gas {
	return co2{}
}

func GetH2OVapour() Gas {
	return h2oVapour{}
}

type air struct{}

func (air) Cp(t float64) float64 {
	// TODO check last value (taken at random)
	var tArr = []float64{
		260, 333, 393, 413, 433, 453, 473, 523, 573, 623, 673, 773, 873, 973, 1073, 1173, 1273, 1373, 1473, 2000,
	}
	var cpArr = []float64{
		1000, 1005, 1009, 1013, 1017, 1022, 1026, 1038, 1047, 1059, 1068, 1093, 1114, 1135, 1156, 1172, 1185, 1197, 1210, 1300,
	}

	var cp = common.InterpTolerate(t, tArr, cpArr)
	return cp
}

func (air) R() float64 {
	return 287.
}

type nitrogen struct{}

func (nitrogen) Cp(t float64) float64 {
	var tArr = []float64{
		275, 300, 325, 350, 375, 400, 450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950, 1000, 1050,
		1100, 1150, 1200, 1250, 1300, 1350, 1400, 1500, 1600, 1700, 1800, 1900, 2000,
	}
	var cpArr = []float64{
		1039, 1040, 1040, 1041, 1042, 1044, 1049, 1056, 1065, 1075, 1086, 1098, 1110, 1122, 1134, 1146,
		1157, 1167, 1177, 1187, 1196, 1204, 1212, 1219, 1226, 1232, 1244, 1254, 1263, 1271, 1278, 1284,
	}

	var cp = common.InterpTolerate(t, tArr, cpArr)
	return cp
}

func (nitrogen) R() float64 {
	return common.UniversalGasConstant / common.N2Weight
}

type co2 struct{}

func (co2) Cp(t float64) float64 {
	var tArr = []float64{
		275, 300, 325, 350, 375, 400, 450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950, 1000, 1050,
		1100, 1150, 1200, 1250, 1300, 1350, 1400, 1500, 1600, 1700, 1800, 1900, 2000,
	}
	var cpArr = []float64{
		819, 846, 871, 895, 918, 939, 978, 1014, 1046, 1075, 1102, 1126, 1148, 1168, 1187, 1204, 1220, 1234,
		1247, 1259, 1270, 1280, 1290, 1298, 1306, 1313, 1326, 1338, 1348, 1356, 1364, 1371,
	}

	var cp = common.InterpTolerate(t, tArr, cpArr)
	return cp
}

func (co2) R() float64 {
	return common.UniversalGasConstant / common.CO2Weight
}

type h2oVapour struct{}

func (h2oVapour) Cp(t float64) float64 {
	var tArr = []float64{
		275, 300, 325, 350, 375, 400, 450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950, 1000, 1050,
		1100, 1150, 1200, 1250, 1300, 1350, 1400, 1500, 1600, 1700, 1800, 1900, 2000,
	}
	var cpArr = []float64{
		819, 846, 871, 895, 918, 939, 978, 1014, 1046, 1075, 1102, 1126, 1148, 1168, 1187, 1204, 1220, 1234,
		1247, 1259, 1270, 1280, 1290, 1298, 1306, 1313, 1326, 1338, 1348, 1356, 1364, 1371,
	}

	var cp = common.InterpTolerate(t, tArr, cpArr)
	return cp
}

func (h2oVapour) R() float64 {
	return common.UniversalGasConstant / common.H2OWeight
}
