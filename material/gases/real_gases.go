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
		260, 333, 393, 413, 433, 453, 473, 523, 573, 623,
		673, 773, 873, 973, 1073, 1173, 1273, 1373, 1473, 2000,
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

func (air) Mu(t float64) float64 {
	var tArr = []float64{
		250, 300,
		350, 400, 450, 500, 550,
		600, 650, 700, 750, 800,
		850, 900, 950, 1000, 1100,
		1200, 1300, 1400, 1500, 1600,
		1700, 1800, 1900, 2000, 2100,
	}
	var muArr = []float64{
		159.6e-7, 184.6e-7,
		208.2e-7, 230.1e-7, 250.7e-7, 270.1e-7, 288.4e-7,
		305.8e-7, 322.5e-7, 338.8e-7, 354.6e-7, 369.8e-7,
		384.3e-7, 398.1e-7, 411.3e-7, 424.4e-7, 449.0e-7,
		473.0e-7, 496.0e-7, 530.0e-7, 557.0e-7, 584.0e-7,
		611e-7, 637e-7, 663e-7, 689e-7, 715e-7,
	}
	return common.InterpTolerate(t, tArr, muArr)
}

func (air) Lambda(t float64) float64 {
	var tArr = []float64{
		250, 300,
		350, 400, 450, 500, 550,
		600, 650, 700, 750, 800,
		850, 900, 950, 1000, 1100,
		1200, 1300, 1400, 1500, 1600,
		1700, 1800, 1900, 2000, 2100,
	}
	var kArr = []float64{
		22.3e-3, 26.3e-3,
		30.0e-3, 33.8e-3, 37.3e-3, 40.7e-3, 43.9e-3,
		46.9e-3, 49.7e-3, 52.4e-3, 54.9e-3, 57.3e-3,
		59.6e-3, 62.0e-3, 64.3e-3, 66.7e-3, 71.5e-3,
		76.3e-3, 82.0e-3, 91.0e-3, 100e-3, 106e-3,
		113e-3, 120e-3, 128e-3, 137e-3, 147e-3,
	}
	return common.InterpTolerate(t, tArr, kArr)
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

func (nitrogen) Mu(t float64) float64 {
	var tArr = []float64{
		300,
		350, 400, 450, 500, 550,
		600, 700, 800, 900, 1000,
		1100, 1200, 1300,
	}
	var muArr = []float64{
		178.2e-7,
		200.0e-7, 220.4e-7, 239.6e-7, 257.7e-7, 274.7e-7,
		290.8e-7, 321.0e-7, 349.1e-7, 375.3e-7, 399.9e-7,
		423.2e-7, 445.3e-7, 466.2e-7,
	}
	return common.InterpTolerate(t, tArr, muArr)
}

func (nitrogen) Lambda(t float64) float64 {
	var tArr = []float64{
		300,
		350, 400, 450, 500, 550,
		600, 700, 800, 900, 1000,
		1100, 1200, 1300,
	}
	var lambdaArr = []float64{
		25.9e-3,
		29.3e-3, 32.7e-3, 35.8e-3, 38.9e-3, 41.7e-3,
		44.6e-3, 49.9e-3, 54.8e-3, 59.7e-3, 64.7e-3,
		70.0e-3, 75.8e-3, 81.0e-3,
	}
	return common.InterpTolerate(t, tArr, lambdaArr)
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

func (co2) Mu(t float64) float64 {
	var tArr = []float64{
		300, 320, 340, 360,
		380, 400, 450, 500, 550,
		600, 650, 700, 750, 800,
	}
	var muArr = []float64{
		149e-7, 156e-7, 165e-7, 173e-7,
		181e-7, 190e-7, 210e-7, 231e-7, 251e-7,
		270e-7, 288e-7, 305e-7, 321e-7, 337e-7,
	}
	return common.InterpTolerate(t, tArr, muArr)
}

func (co2) Lambda(t float64) float64 {
	var tArr = []float64{
		300, 320, 340, 360,
		380, 400, 450, 500, 550,
		600, 650, 700, 750, 800,
	}
	var lambdaArr = []float64{
		16.55e-3, 18.05e-3, 19.70e-3, 21.2e-3,
		22.75e-3, 24.3e-3, 28.3e-3, 32.5e-3, 36.6e-3,
		40.7e-3, 44.5e-3, 48.1e-3, 51.7e-3, 55.1e-3,
	}
	return common.InterpTolerate(t, tArr, lambdaArr)
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

func (h2oVapour) Mu(t float64) float64 {
	var tArr = []float64{
		380, 400, 450, 500, 550,
		600, 650, 700, 750, 800, 850,
	}
	var muArr = []float64{
		127.1e-7, 134.4e-7, 152.5e-7, 170.4e-7, 188.4e-7,
		206.7e-7, 224.7e-7, 242.6e-7, 260.4e-7, 278.6e-7, 296.9e-7,
	}
	return common.InterpTolerate(t, tArr, muArr)
}

func (h2oVapour) Lambda(t float64) float64 {
	var tArr = []float64{
		380, 400, 450, 500, 550,
		600, 650, 700, 750, 800, 850,
	}
	var lambdaArr = []float64{
		24.6e-3, 26.1e-3, 29.9e-3, 33.9e-3, 37.9e-3,
		42.2e-3, 46.4e-3, 50.5e-3, 54.9e-3, 59.2e-3, 63.7e-3,
	}
	return common.InterpTolerate(t, tArr, lambdaArr)
}
