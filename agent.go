package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
)

type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
}

func (a *Agent) Run(ctx context.Context) error {
	log.Println("Running agent")
	conversation := []anthropic.MessageParam{}

	fmt.Println("Chat with claude. Ctrl c quits")

	shouldReadUserInput := true
	for {
		if shouldReadUserInput {
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			input, ok := a.getUserMessage()
			if !ok {
				fmt.Println("Good Bye!")
				break
			}
			userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(input))
			log.Printf("Adding usermessage '%s' to conversation\n", input[:min(len(input), 30)]+"...")
			conversation = append(conversation, userMessage)
		}

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			log.Printf("Inference failed with error: %s\n", err.Error())
			return err
		}

		conversation = append(conversation, message.ToParam())

		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, content := range message.Content {
			log.Printf("Handling content type: %s\n", content.Type)
			switch content.Type {
			case "text":
				fmt.Printf("\u001b[93mClaude\u001b[0m: %s\n", content.Text)
			case "tool_use":
				log.Printf("Should execute tool %s\n", content.Name)
				result := a.executeTool(content.ID, content.Name, content.Input)
				toolResults = append(toolResults, result)
			}
		}
		shouldReadUserInput = len(toolResults) == 0
		if len(toolResults) > 0 {
			log.Printf("Adding tool results to conversation\n")
			conversation = append(conversation, anthropic.NewUserMessage(toolResults...))
		}
	}

	return nil
}

func (a *Agent) executeTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	log.Printf("Searching for a tool with name %s", name)

	tool := a.getToolWithName(name)
	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", name, input)
	response, err := tool.Function(input)
	if err != nil {
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	return anthropic.NewToolResultBlock(id, response, false)
}

func (a *Agent) getToolWithName(name string) ToolDefinition {
	for _, tool := range a.tools {
		if tool.Name == name {
			return tool
		}
	}

	return NoOpTool()
}

func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	return a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Haiku20241022,
		MaxTokens: int64(1024),
		Messages:  conversation,
		Tools:     ToAnthropicTools(a.tools),
	})
}
