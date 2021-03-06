package graph

import (
	"fmt"
)

func Link(port1 Port, port2 Port) {
	port1.SetLinkPort(port2)
	port2.SetLinkPort(port1)

	port1.SetOuterNode(port2.GetInnerNode())
	port2.SetOuterNode(port1.GetInnerNode())
}

func LinkAll(ports1, ports2 []Port) error {
	if len1, len2 := len(ports1), len(ports2); len1 != len2 {
		return fmt.Errorf("length of ports1 %d is not equal to the length of ports2 %d", len1, len2)
	}

	for i := 0; i != len(ports1); i++ {
		Link(ports1[i], ports2[i])
	}
	return nil
}

func SetAll(states []PortState, ports []Port) error {
	if len1, len2 := len(states), len(ports); len1 != len2 {
		return fmt.Errorf("length of states %d is not equal to the length of ports %d", len1, len2)
	}
	for i := 0; i != len(states); i++ {
		ports[i].SetState(states[i])
	}
	return nil
}

func AttachAllWithTags(node Node, ports []*Port, tags []string) {
	AttachAll(node, ports...)
	for i, portPtr := range ports {
		(*portPtr).SetTag(tags[i])
	}
}

func AttachAll(node Node, ports ...*Port) {
	for i := range ports {
		*ports[i] = NewAttachedPort(node)
	}
}

func CopyAll(p1s, p2s []Port) {
	for i, p1 := range p1s {
		CopyState(p1, p2s[i])
	}
}

func CopyState(p1, p2 Port) {
	p2.SetState(p1.GetState())
}

type Port interface {
	GetState() PortState
	SetState(state PortState)
	SetStateNoReverse(state PortState, caller Port)
	GetInnerNode() Node
	SetInnerNode(Node)
	GetOuterNode() Node
	SetOuterNode(Node)
	GetLinkPort() Port
	SetLinkPort(Port)
	GetTag() string
	SetTag(string)
	Attach(another Port)
}

func NewAttachedPortWithTag(node Node, tag string) Port {
	var result = NewAttachedPort(node)
	result.SetTag(tag)
	return result
}

func NewAttachedPort(node Node) Port {
	var port = &portType{
		state:     nil,
		innerNode: nil,
		outerNode: nil,
		linkPort:  nil,
	}
	port.SetInnerNode(node)
	return port
}

func NewPort() Port {
	return &portType{
		state:     nil,
		innerNode: nil,
		outerNode: nil,
		linkPort:  nil,
	}
}

type portType struct {
	state     PortState
	innerNode Node
	outerNode Node
	linkPort  Port
	tag       string

	attachedPorts []Port
}

func (port *portType) Attach(another Port) {
	port.attachedPorts = append(port.attachedPorts, another)
}

func (port *portType) GetTag() string {
	return port.tag
}

func (port *portType) SetTag(tag string) {
	port.tag = tag
}

func (port *portType) GetState() PortState {
	return port.state
}

func (port *portType) SetState(state PortState) {
	port.state = state
	if port.linkPort != nil {
		port.linkPort.SetStateNoReverse(state, port)
	}
	for _, another := range port.attachedPorts {
		another.SetStateNoReverse(state, port)
	}
}

func (port *portType) SetStateNoReverse(state PortState, caller Port) {
	port.state = state
	for _, another := range port.attachedPorts {
		another.SetStateNoReverse(state, caller)
	}
}

func (port *portType) GetInnerNode() Node {
	return port.innerNode
}

func (port *portType) SetInnerNode(src Node) {
	port.innerNode = src
}

func (port *portType) GetOuterNode() Node {
	return port.outerNode
}

func (port *portType) SetOuterNode(dest Node) {
	port.outerNode = dest
}

func (port *portType) GetLinkPort() Port {
	return port.linkPort
}

func (port *portType) SetLinkPort(another Port) {
	port.linkPort = another
}

func (port *portType) setStateNoReverse(state PortState) {
	port.state = state
	for _, another := range port.attachedPorts {
		another.SetState(state)
	}
}
