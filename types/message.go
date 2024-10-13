package types 

import(
	"github.com/google/uuid"
	"context"
)

type Message struct {
	ID uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"user_id"`
	SessionId uuid.UUID `json:"session_id"`
	UserMessage string `json:"user_message"`
	AiMessage string `json:"ai_message"`
}

// for sending messages to the AWS Bedrock InvokeModel endpoint
type MessageRequest struct {
	Text string `json:"inputText"`
}

type MessageList struct {
	Messages []Message `json:"messages"`
}

type MessageTableGateway interface {
	StoreMessage(ctx context.Context, message Message, session Session) (bool, error)
	RetrieveMessages(ctx context.Context, sessionID uuid.UUID) (MessageList, error)
}

type AnthropicMessage struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type AnthropicMessageRequest struct {
	AnthropicVersion string `json:"anthropic_version"`
	MaxTokens int `json:"max_tokens"`
	System string `json:"system"`
	Messages []AnthropicMessage `json:"messages"`
}


type A121Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type A121MessageRequest struct {
	Messages []A121Message `json:"messages"`
	Number int `json:"n"`
}
