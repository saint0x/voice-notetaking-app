package summarization

import (
  "context"
  "log"
  "github.com/sashabaranov/go-openai"
)


// SummarizeText summarizes text using OpenAI API.
func SummarizeText(text string) (string, error) {
  systemPrompt := `
    You are an AI assistant tasked with summarizing a transcription. 
    Give a concise summary of the provided transcription.
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
    return "", err
  }

  return resp.Choices[0].Message.Content, nil
}
