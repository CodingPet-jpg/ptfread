package handler

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Config struct {
	bitmap int64
	sheet  string
	length int
}

var Cfg = Config{}

func init() {
	f, err := os.Open("config.txt")
	if err != nil {
		log.Fatal("Cannot recognize config file : config.txt")
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Println(s)
	}
	f.Close()
}

func (config *Config) HitIndex(i int) bool {
	return (1<<i)&config.bitmap == 1
}

func (config *Config) Sheet() string {
	return config.sheet
}

func (config *Config) Length() int {
	return config.length
}
