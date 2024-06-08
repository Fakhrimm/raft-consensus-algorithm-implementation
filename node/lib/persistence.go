package node

import (
	"Node/grpc/comm"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func (node *Node) LoadLogs(filename string) []comm.Entry {
	log.Printf("Reading saved log")
	filepath := fmt.Sprintf("../config/%vport%v.storage", node.address.IP, node.address.Port)

	var logList []comm.Entry

	file, err := os.OpenFile(filepath, os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to load log file")
		return logList
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	index := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry := node.ParseLog(line)

		logList = append(logList, entry)
		index++
	}

	return logList
}

func (node *Node) SaveLogs() {
	content := ""
	filepath := fmt.Sprintf("../config/%vport%v.storage", node.address.IP, node.address.Port)

	for _, entry := range node.info.log {
		entryString := fmt.Sprintf("%v %v %v\n", entry.Command, entry.Key, entry.Value)
		content += entryString
	}

	contentByte := []byte(content)
	err := os.WriteFile(filepath, contentByte, 0644)

	if err != nil {
		log.Printf("Failed to save log file: %v", err.Error())
		return
	}
}

func (node *Node) SaveLog(entry comm.Entry) {
	content := fmt.Sprintf("%v %v %v\n", entry.Command, entry.Key, entry.Value)
	filepath := fmt.Sprintf("../config/%vport%v.storage", node.address.IP, node.address.Port)

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Printf("Failed to open or create log file: %v", err.Error())
		return
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		log.Printf("Failed to write to file: %v", err.Error())
		return
	}
}

func (node *Node) ParseLog(entryString string) comm.Entry {
	var command int32
	var key string
	var value string
	_, err := fmt.Sscanf(entryString, "%d %s %s", &command, &key, &value)

	if err != nil {
		log.Printf("Error parsing log entry: %v", err)
		return comm.Entry{
			Command: -1,
		}
	}

	return comm.Entry{
		Command: command,
		Key:     key,
		Value:   value,
	}
}
