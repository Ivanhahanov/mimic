package main

import (
	"context"
	"fmt"
	"mimic/pkg/config/mimic"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

func main() {
	var config mimic.Config

	yamlFile, err := os.ReadFile(os.Getenv("CONFIG"))

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, ingress := range config.Mimic.Ingress {
		if ingress.HTTP != nil {
			fmt.Printf("Run TCP server on %d\n", ingress.HTTP.Port)
			wg.Add(1)
			go func() {
				ingress.HTTP.Run()
				wg.Done()
			}()
			// time.Sleep(1 * time.Second)
		} else if ingress.TCP != nil {
			fmt.Printf("Run TCP server on %d\n", ingress.TCP.Port)
		}
	}
	for _, egress := range config.Mimic.Egress {
		if egress.HTTP != nil {
			wg.Add(1)
			go func() {
				egress.HTTP.Run(context.Background())
				wg.Done()
			}()
		} else if egress.TCP != nil {
			fmt.Printf("TCP  request to server on %d\n", egress.TCP.Port)
		}
	}
	wg.Wait()
}
