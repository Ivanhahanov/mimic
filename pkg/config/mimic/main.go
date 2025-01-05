package mimic

import (
	"mimic/pkg/rest"
	"mimic/pkg/ui"
)

type Config struct {
	Mimic Mimic `json:"mimic,omitempty"`
}

type Mimic struct {
	Ingress []Ingress `json:"ingress"`
	Egress  []Egress  `json:"egress"`
}

type Ingress struct {
	HTTP *ui.IngressHTTP `json:"http,omitempty"`
	TCP  *IngressTCP     `json:"tcp,omitempty"`
}

type Egress struct {
	HTTP *rest.Client `json:"http,omitempty"`
	TCP  *EgressTCP   `json:"tcp,omitempty"`
}

type EgressTCP struct {
	Port int `json:"port"`
}

type IngressTCP struct {
	Port int `json:"port"`
}
