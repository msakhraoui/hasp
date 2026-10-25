package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Snortdata struct {
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	Class     string `json:"class"`
	Msg       string `json:"msg"`
	GID       int    `json:"gid"`
	SID       int    `json:"sid"`
	Rev       int    `json:"rev"`
	Priority  int    `json:"priority"`
	SrcAddr   string `json:"src_addr"`
	SrcPort   int    `json:"src_port"`
	DstAddr   string `json:"dst_addr"`
	DstPort   int    `json:"dst_port"`
	Proto     string `json:"proto"`
}

func tailSnort() {
	snortLog, err := os.Open("/var/log/snort/alert_json.txt")
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

	for {
		var snortdata Snortdata
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
