package handler

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"testing"
)

var filter = func(name string) string {
	return name
}

/*
func TestNewFileSet(t *testing.T) {
	fs := NewFileSet("", filter)
	for _, s := range fs {
		fmt.Println(s)
	}
}
*/
var root string = "/home/jokoi/GoProJ"
var testdir = root + "/pkg/linux_amd64/joko.com/demo"

/*func TestPathWalker(t *testing.T) {
	var wf filepath.WalkFunc = func(path string, info fs.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	}
	filepath.Walk(root, wf)
}*/

func TestListDir(t *testing.T) {
	f, _ := os.Open(testdir)
	names, _ := f.Readdirnames(-1)
	fmt.Println(testdir)
	for i, n := range names {
		fmt.Println(i, ":", n)
	}
}

func TestReadDir(t *testing.T) {
	de, _ := os.ReadDir(testdir)
	for _, d := range de {
		fmt.Println(d.Name())
	}
}

func TestReadDirParam(t *testing.T) {
	f, _ := os.Open(root + "/bin")
	de, _ := f.ReadDir(2)
	for _, v := range de {
		fmt.Println(v.Name())
	}
	de2, _ := f.ReadDir(2)
	for _, v := range de2 {
		fmt.Println(v.Name())
	}
}

func TestReadDirFile(t *testing.T) {
	f, err := os.Open("/home/jokoi/GoProJ/pkg/linux_amd64/joko.com/demo/../demo/gomodemo.a")
	fmt.Println(f.Name())
	if err != nil {
		t.Errorf("error:%v\n", err)
	}
	rdf, ok := fs.File(f).(fs.ReadDirFile)
	if ok {
		de, _ := rdf.ReadDir(2)
		for _, d := range de {
			fmt.Println(d)
		}
	}
}

func TestPathClean(t *testing.T) {
	tests := []struct {
		before string
	}{
		{
			before: "/home//jokoi",
		},
		{
			before: "/home/../jokoi",
		},
		{
			before: "/home/./././jokoi",
		},
		{
			before: "/../home/jokoi",
		},
		{
			before: "../home/jokoi",
		},
		{
			before: "",
		},
		{
			before: "/..",
		},
		{
			before: "/.",
		},
	}
	for _, test := range tests {
		after := path.Clean(test.before)
		fmt.Printf("before clean<%s>\tafter clean<%s>\n", test.before, after)
	}
}

func TestModeString(t *testing.T) {
	f, _ := os.Open(testdir)
	fe, _ := f.ReadDir(1)
	s := fe[0].Type().Type()
	fmt.Printf("%d\n", s)
}

func TestDirFS(t *testing.T) {
	de, err := fs.ReadDir(os.DirFS(root), "pkg/linux_amd64")
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
	for _, d := range de {
		fmt.Println(d.Name())
	}
}
