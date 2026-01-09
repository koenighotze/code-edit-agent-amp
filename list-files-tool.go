package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

var ListFilesDefinition = ToolDefinition{
	Name:        "list_files",
	Description: "List all files in a given path below the current working directory.",
	InputSchema: ListFilesInputSchema,
	Function:    ListFiles,
}

type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"the relative path of the directory to list files from in the current working directory. Defaults to the current directory if not specified."`
}

var ListFilesInputSchema = GenerateSchema[ListFilesInput]()

func isDotDirectory(root string, path string, info fs.DirEntry) bool {
	return info.IsDir() && strings.HasPrefix(info.Name(), ".") && path != root
}

func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	err := json.Unmarshal(input, &listFilesInput)
	if err != nil {
		return "", err
	}

	root := "."
	if listFilesInput.Path != "" {
		root = listFilesInput.Path
	}
	log.Printf("Starting walking directory at %s\n", root)
	var files []string
	err = filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		log.Printf("Walking directory %s\n", path)
		if err != nil {
			return err
		}

		if isDotDirectory(root, path, info) {
			log.Printf("Skipping dot directory %s\n", path)
			return filepath.SkipDir
		}

		files = append(files, path)

		return nil
	})

	if err != nil {
		return "", err
	}

	response, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(response), nil
}
