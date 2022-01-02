package handler

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

// provide sync mechanism to call doSimComp

func DoSimComp() {
	var wg sync.WaitGroup
	var fc = make(chan Case, 4000)
	var tokens = make(chan struct{}, 4000)
	var SimCompFunc fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			wg.Add(1)
			tokens <- struct{}{}
			go parse(path, &wg, fc, tokens)
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

		if err := filepath.WalkDir(GetWd(), SimCompFunc); err != nil {
			fmt.Println(err)
		}

		wg.Done()
	}()
	cc := NewCaseChain()
	for done := range fc {
		cc.EliAppend(done)
	}
	for ele := cc.Front(); ele != nil; ele = ele.Next() {
		fmt.Println(ele.Value.(Case).Name)
	}
}

func parse(path string, wg *sync.WaitGroup, parsed chan<- Case, tokens <-chan struct{}) {
	defer func() {
		<-tokens
		wg.Done()
	}()
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := f.GetRows(BaseSheet)
	if err != nil {
		fmt.Println(err)
	}
	var c = NewCase(filepath.Base(path))
	for _, row := range rows {
		if len(row) < Cfg.Length() {
			continue
		}
		var strBuilder bytes.Buffer
		for i, col := range row {
			if Cfg.HitIndex(i) {
				if value, ok := post(col); ok {
					_, err := strBuilder.Write([]byte(value + ","))
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
		// TODO: 当前添加方式可能会有相同条目被添加到手顺节点中，增加之后链表遍历比对的压力，需要在解析时去重
		str := strings.TrimRight(strBuilder.String(), ",")
		if _, ok := c.Contain(str); !ok {
			c.PushBack(str)
		}
	}
	parsed <- c
}
