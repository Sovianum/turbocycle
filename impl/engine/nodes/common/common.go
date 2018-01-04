package common

import "github.com/Sovianum/turbocycle/core/graph"

const (
	DefaultN = 50
)

func IsDataSource(port graph.Port) (bool, error) {
	var linkPort = port.GetLinkPort()
	if linkPort == nil {
		return false, nil
	}

	var outerNode = port.GetOuterNode()
	if outerNode == nil {
		return false, nil
	}

	if !outerNode.ContextDefined() {
		return false, nil
	}

	var updatePorts = outerNode.GetUpdatePorts()

	for _, port := range updatePorts {
		if port == linkPort {
			return true, nil
		}
	}

	return false, nil
}
