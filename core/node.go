package core

type PortsType map[string]Port

type Node interface {
	GetPorts() PortsType
	Process() error //TODO check if need to pass relax coef
	GetRequirePortTags() ([]string,  error)
	GetUpdatePortTags() ([]string, error)
	GetPortTags() []string
	GetPortByTag(tag string) (Port, error)
	ContextDefined() bool
}
