package types 

import (
    "github.com/google/uuid"
	"context"
)

// Metadata represents additional information for an assistant. Used for UI purposes.
type Metadata struct {
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Icon        string    `json:"icon"`
    Prompts     []string  `json:"prompts,omitempty"`
}

// Assistant represents the Bedrock config for an llm assistant.
type Assistant struct {
    ID            uuid.UUID         `json:"id"`
    Name          string            `json:"name"`     // Name of the assistant
    Model         string            `json:"model"`    // Model used by the assistant
    // KbaseID       *uuid.UUID        `json:"kbase_id,omitempty"`
    Type          string            `json:"type"` // Type of the assistant (e.g., travel_assistant, txt-to-sql)
    SystemPrompts string 			`json:"system_prompts"`
    Metadata      *Metadata         `json:"metadata,omitempty"`
}

// AssistantList holds a collection of assistants.
type AssistantList struct {
    Assistants []Assistant `json:"assistants"` // List of assistants
}

type AssistantTableGateway interface {
	CreateAssistant(ctx context.Context, assistant Assistant) (bool, error)
	GetAssistant(ctx context.Context, assistantId uuid.UUID) (Assistant, error)
	UpdateAssistant(ctx context.Context, assistant Assistant) (bool, error)
    ListAssistants(ctx context.Context) (AssistantList, error)
	// GetAssistantByName(ctx context.Context, assistantName string) (Session, error)
}