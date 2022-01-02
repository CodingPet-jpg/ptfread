package handler

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Sheet   string `yaml:"ActiveSheet"`
	Bitmap  int64  `yaml:"BitMap"`
	Length  int    `yaml:"Length"`
	WorkDir string `yaml:"WorkDirectory"`
	Pn      int    `yaml:"ParallelNum"`
}

var Cfg Config

func init() {
	yml, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println(err)
	}
	if err := yaml.Unmarshal(yml, &Cfg); err != nil {
		fmt.Println(err)
	}
}

func (config *Config) HitIndex(i int) bool {
	return (1<<i)&config.Bitmap == 1
}

// defined as variable for easy switch when doing test

type PostFunc func(name string) (string, bool)

var post PostFunc = func(name string) (string, bool) {
	return name, true
}

func ChangePoster(f PostFunc) {
	post = f
}
