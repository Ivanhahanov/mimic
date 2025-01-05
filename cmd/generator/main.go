package main

import (
	"fmt"
	config "mimic/pkg/config/generator"
	"mimic/pkg/generator"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	var config config.Generator

	yamlFile, err := os.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	fmt.Println("%v", config)
	for _, ns := range config.Namespaces {
		for _, resource := range ns.Resources {
			if resource.Deployment != nil {
				manifest := resource.Deployment.Generate(ns.Name, config.Image)
				generator.Save(ns.Name, "deployment", resource.Deployment.Name, manifest)
				cm := resource.Deployment.GenerateConfigMap(ns.Name)
				generator.Save(ns.Name, "cm", resource.Deployment.Name, cm)
				svc := resource.Deployment.GenerateService(ns.Name)
				generator.Save(ns.Name, "svc", resource.Deployment.Name, svc)
			}
		}
	}
}
