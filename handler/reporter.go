package handler

import (
	"bytes"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"text/tabwriter"
	"time"
)

type conclusion struct {
	cost        time.Duration
	remainCount int
}

var brief conclusion
var start time.Time

func report(chain CaseChain) conclusion {
	var location = filepath.Join(Cfg.WorkDir, "report")
	brief.remainCount = chain.Len()
	errm := os.Mkdir(location, os.ModeDir)
	if !errors.Is(errm, fs.ErrExist) {
		log.Printf("Failed to create directory %s", location)
	}
	var filename = "report" + "[" + time.Now().Format("2006-01-02 ※ 15-04-05") + "]" + ".txt"
	f, erro := os.OpenFile(filepath.Join(location, filename), os.O_RDWR|os.O_CREATE, 0666)
	if erro != nil {
		log.Printf("Failed to create report:%s\n", filename)
	}
	tw := new(tabwriter.Writer).Init(f, 0, 8, 2, ' ', 0)
	bufc := make(chan *bytes.Buffer, 100)
	wg := sync.WaitGroup{}
	for ele := chain.Front(); ele != nil; ele = ele.Next() {
		wg.Add(1)
		go ele.Value.(Case).Format(f, bufc, &wg)
	}
	go func() {
		wg.Wait()
		close(bufc)
	}()
	for buf := range bufc {
		tw.Write([]byte(buf.String()))
	}
	log.Printf("Detailed report located at %s\n", filepath.Join(location, filename))
	return brief
}

func begin() {
	log.Printf("Start in work directory %s\n", Cfg.WorkDir)
	start = time.Now()
}

func end() {
	log.Println("All task completed,Taking a while to generate the report")
	brief.cost = time.Since(start)
}
