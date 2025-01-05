package generator

import "mimic/pkg/generator"

type Generator struct {
	Image      string      `yaml:"image,omitempty"`
	Namespaces []Namespace `yaml:"namespaces,omitempty"`
}

type Namespace struct {
	Name      string     `yaml:"name,omitempty"`
	Resources []Resource `yaml:"resources,omitempty"`
}

type Resource struct {
	Deployment *generator.Deployment `yaml:"deployment"`
}
