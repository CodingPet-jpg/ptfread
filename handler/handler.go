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

func ParllelComp() {
	var (
		wg     sync.WaitGroup
		parsed = make(chan Case, Cfg.Pn)
	)

	var parse = func(path string) {
		defer func() {
			wg.Done()
		}()
		f, err := excelize.OpenFile(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows, err := f.GetRows(Cfg.Sheet)
		if err != nil {
			fmt.Println(err)
		}
		var tcase = NewCase(filepath.Base(path))
		for _, row := range rows {
			if len(row) < Cfg.Length {
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
			str := strings.TrimRight(strBuilder.String(), ",")
			if _, ok := tcase.Contain(str); !ok {
				tcase.PushBack(str)
			}
		}
		parsed <- tcase
	}

	var SimCompFunc fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !strings.HasPrefix(d.Name(), "~") && (strings.HasSuffix(d.Name(), ".xlsm") || strings.HasSuffix(d.Name(), ".xlsx")) {
			wg.Add(1)
			go parse(path)
		}
		return nil
	}

	go func() {
		wg.Add(1)

		go func() {
			wg.Wait()
			close(parsed)
		}()

		if err := filepath.WalkDir(Cfg.WorkDir, SimCompFunc); err != nil {
			fmt.Println(err)
		}

		wg.Done()
	}()

	cc := NewCaseChain()
	for done := range parsed {
		cc.EliAppend(done)
	}
	for ele := cc.Front(); ele != nil; ele = ele.Next() {
		fmt.Println(ele.Value.(Case).Name)
	}
}