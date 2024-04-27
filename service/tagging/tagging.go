package tagging

import (
	"context"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"strings"
)

// TagText tags text with topics using OpenAI API.
func TagText(text string) ([]string, error) {
	systemPrompt := `
      You are an AI assistant tasked with extracting tags or topics from a transcription. 
      List the relevant tags or topics based on the provided transcription.
    `

	client := openai.NewClient("your_token_here")

	// Prepare messages
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		},
	}

	// Create Chat completion request
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return nil, err
	}

	// Extracting tags from the response
	rawTags := resp.Choices[0].Message.Content
	tags := strings.Split(rawTags, "\n")

	return tags, nil
}
