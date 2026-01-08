package main

import (
	"encoding/json"
	"log"
	"os"
)

var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"the relative path of the file to read in the current working directory."`
}

var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

func ReadFile(input json.RawMessage) (string, error) {
	rfInput := ReadFileInput{}

	err := json.Unmarshal(input, &rfInput)
	if err != nil {
		panic(err)
	}

	log.Printf("Should read file %s for agent\n", rfInput.Path)
	content, err := os.ReadFile(rfInput.Path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
