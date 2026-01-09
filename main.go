package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	logFile, err := os.OpenFile("agent.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting anthropic agent poc")
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)

	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	tools := []ToolDefinition{ReadFileDefinition, ListFilesDefinition}
	agent := NewAgent(&client, getUserMessage, tools)

	err = agent.Run(context.TODO())
	if err != nil {
		panic(err)
	}
}

func NewAgent(client *anthropic.Client, getUserMessage func() (string, bool), tools []ToolDefinition) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}
