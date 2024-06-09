package node

import (
	"Node/grpc/comm"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func (node *Node) LoadLogs() {
	log.Printf("Reading saved log")
	filepath := fmt.Sprintf("../config/%vport%v.storage", node.address.IP, node.address.Port)

	var logList []comm.Entry

	file, err := os.OpenFile(filepath, os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to load log file")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	index := 0
	newConfigIdx := -1
	oldnewConfigIdx := -1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry := node.ParseLog(line)

		if entry.Command == int32(NewOldConfig) {
			oldnewConfigIdx = index
		} else if entry.Command == int32(NewConfig) {
			newConfigIdx = index
		}

		logList = append(logList, entry)
		index++
	}

	if oldnewConfigIdx > newConfigIdx {
		node.info.isJointConsensus = true
		node.info.newClusterAddresses = parseAddresses(logList[oldnewConfigIdx].Value)
		node.info.newClusterCount = len(node.info.newClusterAddresses)
		node.info.newId = -1
		for index, address := range node.info.newClusterAddresses {
			if node.address.String() == address.String() {
				node.info.newId = index
				break
			}
		}
	} else if newConfigIdx > oldnewConfigIdx && oldnewConfigIdx != -1 {
		node.info.newClusterAddresses = parseAddresses(logList[oldnewConfigIdx].Value)
		node.info.newClusterCount = len(node.info.newClusterAddresses)
		node.info.newId = -1
		for index, address := range node.info.newClusterAddresses {
			if node.address.String() == address.String() {
				node.info.newId = index
				break
			}
		}

		node.info.isJointConsensus = false
		node.info.clusterCount = node.info.newClusterCount
		node.info.clusterAddresses = node.info.newClusterAddresses

		if node.info.newId == -1 {
			defer node.Stop()
		} else {
			node.info.id = node.info.newId
		}
	}
	node.info.log = logList

	// Set current term to the last term in the log
	node.info.currentTerm = int(node.info.log[len(node.info.log)-1].Term)
}

func (node *Node) SaveLogs() {
	content := ""
	filepath := fmt.Sprintf("../config/%vport%v.storage", node.address.IP, node.address.Port)

	for _, entry := range node.info.log {
		entryString := fmt.Sprintf("%v %v %v %v\n", entry.Command, entry.Key, entry.Value, entry.Term)
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
	// TODO: Should config commands be kept persistent?

	// if entry.Command == int32(NewConfig) || entry.Command == int32(NewOldConfig) {
	// 	log.Printf("[Persistence] Log not written because it's a command")
	// 	return
	// }

	content := fmt.Sprintf("%v %v %v %v\n", entry.Command, entry.Key, entry.Value, entry.Term)
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
	var term int32
	_, err := fmt.Sscanf(entryString, "%d %s %s %d", &command, &key, &value, &term)

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
		Term:    term,
	}
}
