package handler

import (
	"fmt"
	"io/fs"
	"log"
	"sync"

	"github.com/CodingPet-jpg/ptfread/base"
	"github.com/xuri/excelize/v2"
)

// provide sync mechanism to call doSimComp
func DoSimComp() {
	var wg sync.WaitGroup
	var fc = make(chan base.Case, 100)
	var tokens = make(chan struct{}, 40)
	var SimCompFunc fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			wg.Add(1)
			tokens <- struct{}{}
			go doFileParse(path, &wg, fc, tokens)
		}
		return nil
	}

	go func() {
		wg.Add(1)

		go func() {
			wg.Wait()
			close(fc)
			close(tokens)
		}()

		fs.WalkDir(base.Fs, ".", SimCompFunc)

		wg.Done()
	}()

	for done := range fc {
		log.Println(done)
	}
}

func doFileParse(path string, wg *sync.WaitGroup, parsed chan<- base.Case, tokens <-chan struct{}) {
	defer func() {
		<-tokens
		wg.Done()
	}()
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := f.Rows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
		}
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
	if err = rows.Close(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(f)
}

/*
type ResultSet struct{}

func NewFileSet(dir string, filter func(filename string) string) []string {
	if !filepath.IsAbs(dir) {
		prefix, _ := os.Getwd()
		dir = prefix + dir
	}
	if fileinfo, err := os.Stat(dir); err != nil {
		log.Fatalf("Failed:%v", err)
	} else if !fileinfo.IsDir() {
		log.Fatalln("Not Dictionary")
	}
	fc := make(chan string, 100)
	wg.Add(1)
	go getFileSet(dir, fc)

	var rs []string
	go func() {
		wg.Wait()
		close(fc)
	}()
	for s := range fc {
		s = filter(s)
		if s != "" {
			rs = append(rs, "/"+s)
		}
	}

	return rs
}

func getFileSet(dir string, fc chan<- string) {
	defer wg.Done()
	direntries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("error:%v", err)
	}
	for _, de := range direntries {
		if de.IsDir() {
			wg.Add(1)
			go getFileSet(dir+de.Name(), fc)
		}
		fc <- dir + de.Name()
	}
}
*/
