package handler

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/xuri/excelize/v2"
)

// provide sync mechanism to call doSimComp

func ParllelComp() {
	begin()
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
		// init the new case with file base name
		var tcase = NewCase(filepath.Base(path))
		for _, row := range rows {
			if len(row) < Cfg.Length {
				continue
			}
			var entry = make([]string, 0, 4)
			for i, col := range row {
				// append target column each row into string slice
				if Cfg.HitIndex(uint64(i)) {
					if value, ok := post(col); ok {
						entry = append(entry, value)
					}
				}
			}
			if _, ok := tcase.Contain(entry); !ok {
				tcase.PushBack(entry)
			}
		}
		parsed <- tcase
	}

	var SimCompFunc fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !strings.HasPrefix(d.Name(), "~") && (strings.HasSuffix(d.Name(), ".xlsm") || strings.HasSuffix(d.Name(), ".xlsx")) {
			wg.Add(1)
			atomic.AddInt64(&fileCounter, 1)
			go parse(path)
		}
		return nil
	}

	if Cfg.WorkDir != "" {
		wg.Add(1)
		go func() {
			if err := filepath.WalkDir(Cfg.WorkDir, SimCompFunc); err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}()
	}

	var chain = NewCaseChain()
	// TODO:support incremental updating
	wg2 := sync.WaitGroup{}
	var parsed2 chan Case

	if len(Cfg.Is) != 0 {
		parsed2 = make(chan Case, len(Cfg.Is))
		for i, p := range Cfg.Is {
			f, err := os.Open(p)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if i == 0 {
				wg2.Add(1)
				go func() {
					wg2.Wait()
					close(parsed2)
				}()
				go func(file *os.File) {
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						ucase := UnMarshal([]byte(scanner.Text()))
						parsed2 <- ucase
					}
					wg2.Done()
				}(f)
			} else {
				wg.Add(1)
				go func(file *os.File) {
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						ucase := UnMarshal([]byte(scanner.Text()))
						parsed <- ucase
					}
					wg.Done()
				}(f)
			}
		}
	}
	go func() {
		wg.Wait()
		close(parsed)
	}()
	if parsed2 != nil {
		for pcase2 := range parsed2 {
			chain.PushBack(pcase2)
		}
	}
	// only WalkDir can release parsed channel
	for pcase := range parsed {
		chain.EliAppend(pcase)
	}

	report(chain)
}
