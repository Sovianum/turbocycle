package core

import "errors"

type Network struct {
	nodes []Node
}

func (network *Network) checkFreePorts() error {
	for _, node := range network.nodes {
		for _, port := range node.GetPorts() {
			if port.linkPort == nil {
				return errors.New("Found free port")
			}
		}
	}

	return nil
}
