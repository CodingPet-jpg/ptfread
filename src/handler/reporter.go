package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var start time.Time
var fileCounter int64 = 0
var wd string

func report(chain CaseChain) {
	log.Printf("All task completed\n  <Handled File Count:%d>\n  <Remaining File Count:%d>\n  <Spend Time:%s>\nTaking a while to generate the report\n", fileCounter, chain.Len(), time.Since(start))
	var (
		location         = filepath.Join(wd, "report")
		filename         = "report" + "[" + time.Now().Format("2006-01-02 ※ 15-04-05") + "]" + ".txt"
		absoluteFilePath = filepath.Join(location, filename)
	)
	errm := os.Mkdir(location, os.ModeDir)
	if errm != nil && !errors.Is(errm, fs.ErrExist) {
		log.Printf("Failed to create directory %s", location)
	}

	f, erro := os.OpenFile(absoluteFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if erro != nil {
		log.Printf("Failed to create report:%s\n", filename)
	}
	defer f.Close()
	// the size of channel use half of parallel number
	bufc := make(chan *bytes.Buffer, Cfg.Pn/2)
	wg := sync.WaitGroup{}
	for ele := chain.Front(); ele != nil; ele = ele.Next() {
		wg.Add(1)
		go ele.Value.(Case).Marshal(bufc, &wg)
	}
	go func() {
		wg.Wait()
		close(bufc)
	}()

	for buf := range bufc {
		_, err := f.Write([]byte(buf.String()))
		if err != nil {
			fmt.Println(err)
		}
	}

	log.Printf("Detailed report located at %s\n", absoluteFilePath)
}

func begin() {
	if Cfg.WorkDir == "" {
		wd, _ = os.Getwd()
	} else {
		wd = Cfg.WorkDir
	}
	log.Printf("Starting in work directory %s\n", wd)
	start = time.Now()
}
