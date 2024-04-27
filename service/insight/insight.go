package insight

import (
  "context"
  "fmt"
  "github.com/sashabaranov/go-openai"
  "log"
)


// GenerateInsight generates insights from grouped notes using OpenAI API.
func GenerateInsight(groupedNotes []string) (string, error) {
  systemPrompt := `
    You are an AI assistant tasked with generating insights from grouped notes. 
    Provide insights based on the provided grouped notes.
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
      Content: strings.Join(groupedNotes, "\n"),
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

// GenerateInsightByTime generates insights from notes within a specified timeframe using OpenAI API.
func GenerateInsightByTime(notes []string, startTime, endTime string) (string, error) {
  client := openai.NewClient("your_token_here")

  // Prepare messages
  messages := []openai.ChatCompletionMessage{}
  for _, note := range notes {
    messages = append(messages, openai.ChatCompletionMessage{
      Role:    openai.ChatMessageRoleUser,
      Content: note,
    })
  }
  timeframe := fmt.Sprintf("Start Time: %s\nEnd Time: %s", startTime, endTime)
  messages = append(messages, openai.ChatCompletionMessage{
    Role:    openai.ChatMessageRoleUser,
    Content: timeframe,
  })

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
