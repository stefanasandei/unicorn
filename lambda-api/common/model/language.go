package model

import (
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v2"
)

type Language struct {
	Name      string    `yaml:"name"`
	Versions  []float32 `yaml:"versions"`
	Extension string    `yaml:"extension"`

	CompileCmds []string `yaml:"compile,omitempty"`
	RunCmds     []string `yaml:"run"`
}

func NewLanguage(runtimeFile fs.DirEntry, dir string) (Language, error) {
	filepath := fmt.Sprintf("%s/%s", dir, runtimeFile.Name())

	data, err := os.ReadFile(filepath)
	if err != nil {
		return Language{}, err
	}

	var lang Language
	err = yaml.Unmarshal(data, &lang)
	if err != nil {
		return Language{}, err
	}

	return lang, nil
}

/*
	Two special keywords in the compile & run commands arrays:
		- <entry>: the main source file (usually ./worker/tasks/exec_id/main)
		- <output>: the compiled binary after the compilation phase
*/
