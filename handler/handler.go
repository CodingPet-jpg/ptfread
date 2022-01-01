package handler

import (
	"bytes"
	"fmt"
	"io/fs"
	syspath "path"
	"sync"

	"github.com/xuri/excelize/v2"
)

// provide sync mechanism to call doSimComp
func DoSimComp() {
	var wg sync.WaitGroup
	var fc = make(chan *Case, 40)
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

		fs.WalkDir(Fs, ".", SimCompFunc)

		wg.Done()
	}()
	ln := &LinkedNode{}
	for done := range fc {
		ln.ComparedAppend(done)
	}
}

func doFileParse(path string, wg *sync.WaitGroup, parsed chan<- *Case, tokens <-chan struct{}) {
	defer func() {
		<-tokens
		wg.Done()
	}()
	f, err := excelize.OpenFile(syspath.Join(GetWd(), path))
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := f.GetRows(BaseSheet)
	if err != nil {
		fmt.Println(err)
	}
	var c = &Case{Name: path}
	for _, row := range rows {
		if len(row) < Cfg.Length() {
			continue
		}
		var strBuilder bytes.Buffer
		for i, col := range row {
			if Cfg.HitIndex(i) {
				if value, ok := regex(col); ok {
					_, err := strBuilder.Write([]byte(value))
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
		// TODO: 当前添加方式可能会有相同条目被添加到手顺节点中，增加之后链表遍历比对的压力，需要在解析时去重
		str := strBuilder.String()
		if _, ok := c.Contain(str); !ok {
			c.PushBack(strBuilder.String())
		}
	}
	parsed <- c
}

func regex(name string) (string, bool) {
	return name, true
}
