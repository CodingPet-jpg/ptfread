package handler

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
)

type Config struct {
	bitmap int64
	sheet  string
	length int
}

// todo:for test ,right code is Config{}
var Cfg = Config{length: 6}

var wd = flag.String("w", "../../testdata", "Set WorkDirectory")

var BaseSheet string

func init() {
	testing.Init()
	if !flag.Parsed() {
		flag.Parse()
	}
	BaseSheet = "Sheet1"
}

func init() {
	f, err := os.Open("./config.txt")
	if err != nil {
		log.Fatal("Cannot recognize config file : config.txt")
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Println(s)
	}
	if err := f.Close(); err != nil {
		fmt.Println(err)
	}

}

func (config *Config) HitIndex(i int) bool {
	//todo:for test,should change to 1
	return (1<<i)&config.bitmap == 0
}

func (config *Config) Sheet() string {
	return config.sheet
}

func (config *Config) Length() int {
	return config.length
}

func GetWd() string {
	return *wd
}

// defined as variable for easy switch when doing test

type PostFunc func(name string) (string, bool)

var post PostFunc = func(name string) (string, bool) {
	return name, true
}

func ChangePoster(f PostFunc) {
	post = f
}
