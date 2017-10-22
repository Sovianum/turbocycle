package nodes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
)

type TurbineStageNode interface {
	core.Node
	nodes.GasChannel
}
