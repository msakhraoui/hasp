package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type CowrieLogs struct {
	EventID   string `json:"eventid"`
	SrcIP     string `json:"src_ip"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Input     string `json:"input"`
	Timestamp string `json:"timestamp"`
}

func tailCowrie() {
	cowrieLog, err := os.Open("/home/cowrie/cowrie/var/log/cowrie/cowrie.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := cowrieLog.Close(); err != nil {
			fmt.Println("error closing file", err)
		}
	}()
	reader := bufio.NewReader(cowrieLog)
	var cowrieLogs CowrieLogs
	for {
		line, no := reader.ReadString('\n')
		parseErr := json.Unmarshal([]byte(line), &cowrieLogs)
		if no != nil {
			time.Sleep(500 * time.Millisecond)
		}
		if parseErr == nil {
			fmt.Println(cowrieLogs)
		}
	}
}
