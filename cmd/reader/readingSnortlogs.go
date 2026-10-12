package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func readingSnortlogs() {
	type Snortdata struct {
	}

	snortLog, err := os.Open("/var/log/snort/alert-fast.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer snortLog.Close()
	reader := bufio.NewReader(snortLog)
	var snortdata Snortdata

	for {
		line, err := reader.ReadString('\n')
		parsegoterr := json.Unmarshal([]byte(line), &snortdata)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
		}
		if parsegoterr == nil { //parsing went good
			fmt.Println(snortdata)
		}
	}
}
