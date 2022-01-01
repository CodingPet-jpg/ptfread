package base

import (
	"flag"
	"io/fs"
	"os"
)

var wd = flag.String("w", ".", "Set WorkDirectory")

var Fs fs.FS

type Action struct {
	Application string
	UI          string
	Operand     string
	Action      string
	Operate     string
}

type Case map[string][]Action

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
	Fs = os.DirFS(*wd)
}

func GetWd() string {
	return *wd
}
