package dotconfig_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goaux/dotconfig"
)

func ExampleDir() {
	dir, exists := dotconfig.Dir("myapp")
	if !exists {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}
	fmt.Println(dir)
}

func ExampleFile() {
	file, status := dotconfig.File("myapp", "config.yaml")
	if status == dotconfig.NotExists {
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			panic(err)
		}
	}
	var config []byte
	if err := os.WriteFile(file, config, 0644); err != nil {
		panic(err)
	}
	fmt.Println(file)
}
