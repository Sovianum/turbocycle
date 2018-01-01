package core

type PortsType map[string]Port

type Node interface {
	//GetPorts() PortsType
	Process() error
	GetRequirePorts() []Port
	GetUpdatePorts() []Port
	GetPorts() []Port
	//GetRequirePortTags() ([]string, error)
	//GetUpdatePortTags() ([]string, error)
	//GetPortTags() []string
	GetPortByTag(tag string) (Port, error)
	ContextDefined() bool
}
