package definitions

import (
	"smr/pkg/network"
)

// Local contracts
type Containers struct {
	Kind       string               `yaml:"kind"`
	Containers map[string]Container `mapstructure:"container"`
}

type Container struct {
	Meta Meta `yaml:"meta"`
	Spec Spec `yaml:"spec"`
}

type Meta struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
	Group   string `yaml:"group"`
}

type Spec struct {
	Container ContainerInternal `mapstructure:"containerx"`
}

type ContainerInternal struct {
	Image         string                 `yaml:"image"`
	Tag           string                 `yaml:"tag""`
	Envs          []string               `yaml:"envs"`
	Entrypoint    []string               `yaml:"entrypoint"`
	Cmd           []string               `yaml:"cmd"`
	Dependencies  []DependsOn            `yaml:"dependencies"`
	Networks      []string               `yaml:"networks"`
	Ports         []network.PortMappings `yaml:"ports"`
	FileMounts    []map[string]string    `yaml:"fileMounts"`
	Operators     []map[string]any       `yaml:"operators"`
	Configuration map[string]any         `mapstructure:"configuration"`
	Resources     []map[string]any       `mapstructure:"resources"`
	Replicas      int                    `yaml:"replicas"`
}

type DependsOn struct {
	Name     string         `yaml:"name"`
	Operator string         `yaml:"operator"`
	Timeout  string         `yaml:"timeout"`
	Body     map[string]any `mapstructure:"body"`
	Solved   bool
}

// Internal implementation contracts

// Configuration

type Configuration struct {
	Meta ConfigurationMeta `mapstructure:"meta"`
	Spec ConfigurationSpec `mapstructure:"spec"`
}

type ConfigurationMeta struct {
	Group      string `json:"group"`
	Identifier string `json:"identifier"`
}

type ConfigurationSpec struct {
	Data map[string]string `json:"data"`
}

// Template

type Resource struct {
	Meta ResourceMeta `mapstructure:"meta"`
	Spec ResourceSpec `mapstructure:"spec"`
}

type ResourceMeta struct {
	Group      string `json:"group"`
	Identifier string `json:"identifier"`
}

type ResourceSpec struct {
	Data map[string]string `json:"data"`
}
