package common

import "github.com/Sovianum/turbocycle/core/graph"

type VelocityChannel interface {
	VelocitySink
	VelocitySource
}

type VelocitySink interface {
	VelocityInput() graph.Port
}

type VelocitySource interface {
	VelocityOutput() graph.Port
}

type DimensionChannel interface {
	DimensionSource
	DimensionSink
}

type DimensionSource interface {
	DimensionInput() graph.Port
}

type DimensionSink interface {
	DimensionSource() graph.Port
}
