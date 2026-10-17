package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {
	type Snortdata struct {
		Timestamp string `json:"timestamp"`
		Src_addr  string `json:"src_addr"`
		Src_port  string `json:"src_port"`
		Dst_addr  string `json:"dst_addr"`
		Dst_port  string `json:"dst_port"`
		Rule      string `json:"rule"`
		Sid       string `json:"sid"`
		Action    string `json:"action"`
		Msg       string `json:"msg"`
	}

	snortLog, err := os.Open("/var/log/snort/alert-fast.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := snortLog.Close(); err != nil {
			fmt.Println("error closing file", err)
		}
	}()

	reader := bufio.NewReader(snortLog)
	var snortdata Snortdata

	for {
		line, err := reader.ReadString('\n')
		parsegoterr := json.Unmarshal([]byte(line), &snortdata)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
		}
		if parsegoterr == nil { // parsing went good
			fmt.Println(snortdata)
		}
	}
}
