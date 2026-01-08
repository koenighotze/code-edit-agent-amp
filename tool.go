package main

import (
	"encoding/json"
	"errors"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

func NoOpTool() ToolDefinition {
	return ToolDefinition{
		Name:        "no_op",
		Description: "does nothing",
		Function:    func(input json.RawMessage) (string, error) { return "nothing", errors.New("tool not found") },
	}
}

func ToAnthropicTools(tools []ToolDefinition) []anthropic.ToolUnionParam {
	aTools := []anthropic.ToolUnionParam{}
	for _, tool := range tools {
		aTools = append(aTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	return aTools
}

func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var v T
	schema := reflector.Reflect(v)
	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}
