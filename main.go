package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	interval = time.Second * 5
	write    = ""
	read     = ""

	// Channel to shut down
	done  = make(chan bool)
	start = time.Now()
)

func init() {
	flag.StringVarP(&write, "write", "w", "", "Write to this file")
	flag.StringVarP(&read, "read", "r", "", "Read from this file")
	flag.DurationVarP(&interval, "interval", "i", 5*time.Second, "Interval to read/write")
	flag.Parse()
}

func main() {

	if write == "" && read == "" {
		fmt.Println("You must specify --write and/or --read")
		flag.Usage()
		os.Exit(1)
	}

	if write != "" {
		fmt.Printf("Writing every %v\n", interval)
	}

	if read != "" {
		fmt.Printf("Reading every %v\n", interval)
	}

	// Catch ctrl-c / sigterm
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt
		done <- true
	}()

	loop()

	fmt.Println("Done.")
}

func loop() {
	doReadWrite()

	// Read/write every interval
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			doReadWrite()
		}
	}
}

func doReadWrite() {
	if write != "" {
		doWrite()
	}
	if read != "" {
		doRead()
	}
}

func doRead() {
	data, err := ioutil.ReadFile(read)
	if err != nil {
		panic(err)
	}

	timeStr := strings.TrimSpace(string(data))
	timeTime, err := time.Parse(time.RFC3339Nano, timeStr)
	since := time.Since(timeTime)
	if err != nil {
		panic(err)
	}

	if since <= interval {
		fmt.Printf("File contents OK (%v old)\n", since)
	} else {
		fmt.Printf("File contents BAD (%v old, %v since start)\n", since, time.Since(start))
	}
}

func doWrite() {
	timeStr := time.Now().UTC().Format(time.RFC3339Nano)
	err := ioutil.WriteFile(write, []byte(timeStr+"\n"), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wrote data: %v\n", timeStr)
}
