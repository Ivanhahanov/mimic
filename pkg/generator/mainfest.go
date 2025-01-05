package generator

import (
	"fmt"
	"os"
	"path"
)

func Save(ns string, resource string, name string, manifest []byte) {
	dir := "manifests"
	os.MkdirAll(path.Join(dir, ns), os.ModePerm)
	filename := path.Join(dir, ns, fmt.Sprintf("%s-%s.yml", resource, name))
	err := os.WriteFile(filename, manifest, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
