package dataframes

type ThreeShaftsDF struct {
	GasSource           GasDF
	InletFilter         PressureDropDF
	InletPipe PressureDropDF

	LPCompressor     CompressorDF
	LPCompressorPipe PressureDropDF
	LPTurbine        BlockedTurbineDF
	LPTurbinePipe    PressureDropDF
	LPShaft          ShaftDF

	HPCompressor  CompressorDF
	HPTurbine     BlockedTurbineDF
	HPTurbinePipe PressureDropDF
	HPShaft       ShaftDF

	Burner        BurnerDF

	FreeTurbine FreeTurbineDF
	OutletPipe  PressureDropDF
}
