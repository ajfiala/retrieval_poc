package db

import (
	"context"
	"rag-demo/types"
	"fmt"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

// UserTableGatewayImpl is the implementation of UserTableGateway using pgxpool.
type AssistantTableGatewayImpl struct {
	Pool *pgxpool.Pool
}

// NewUserTableGateway creates a new instance of UserTableGatewayImpl.
func NewAssistantTableGateway(pool *pgxpool.Pool) types.AssistantTableGateway {
	return &AssistantTableGatewayImpl{Pool: pool}
}

func (atg *AssistantTableGatewayImpl) CreateAssistant(ctx context.Context, assistant types.Assistant) (bool, error) {
    systemPromptsJSON, err := json.Marshal(assistant.SystemPrompts)
    if err != nil {
        return false, fmt.Errorf("failed to marshal SystemPrompts: %v", err)
    }

    metadataJSON, err := json.Marshal(assistant.Metadata)
    if err != nil {
        return false, fmt.Errorf("failed to marshal Metadata: %v", err)
    }

    // Execute the SQL query with marshaled JSON
    _, err = atg.Pool.Exec(ctx,
        `INSERT INTO assistant (uuid, name, model, type, system_prompts, metadata)
         VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb)`,
        assistant.ID,  assistant.Name, assistant.Model, assistant.Type, systemPromptsJSON, metadataJSON)
    if err != nil {
        return false, err
    }
    return true, nil
}

func (atg *AssistantTableGatewayImpl) GetAssistant(ctx context.Context, assistantId uuid.UUID) (types.Assistant, error) {
	var assistant types.Assistant
	var systemPrompts string
	var metadata string
	err := atg.Pool.QueryRow(ctx, "SELECT uuid, name, model, type, system_prompts, metadata FROM assistant WHERE uuid = $1", assistantId).Scan(&assistant.ID, &assistant.Name, &assistant.Model, &assistant.Type, &systemPrompts, &metadata)
	if err != nil {
		return types.Assistant{}, err
	}

	err = json.Unmarshal([]byte(systemPrompts), &assistant.SystemPrompts)
	if err != nil {
		return types.Assistant{}, fmt.Errorf("failed to unmarshal SystemPrompts: %v", err)
	}

	err = json.Unmarshal([]byte(metadata), &assistant.Metadata)
	if err != nil {
		return types.Assistant{}, fmt.Errorf("failed to unmarshal Metadata: %v", err)
	}

	return assistant, nil
}

func (atg *AssistantTableGatewayImpl) UpdateAssistant(ctx context.Context, assistant types.Assistant) (bool, error) {
	systemPromptsJSON, err := json.Marshal(assistant.SystemPrompts)
	if err != nil {
		return false, fmt.Errorf("failed to marshal SystemPrompts: %v", err)
	}

	metadataJSON, err := json.Marshal(assistant.Metadata)
	if err != nil {
		return false, fmt.Errorf("failed to marshal Metadata: %v", err)
	}

	// Execute the SQL query with marshaled JSON
	_, err = atg.Pool.Exec(ctx,
		`UPDATE assistant SET name = $2, model = $3, type = $4, system_prompts = $5::jsonb, metadata = $6::jsonb WHERE uuid = $1`,
		assistant.ID,  assistant.Name, assistant.Model, assistant.Type, systemPromptsJSON, metadataJSON)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (atg *AssistantTableGatewayImpl) ListAssistants(ctx context.Context) (types.AssistantList, error) {
    rows, err := atg.Pool.Query(ctx, "SELECT uuid, name, model, type, system_prompts, metadata FROM assistant")
    if err != nil {
        return types.AssistantList{}, err
    }
    defer rows.Close()

    var assistants []types.Assistant
    for rows.Next() {
        var assistant types.Assistant
        var systemPrompts string
        var metadata string

        err := rows.Scan(&assistant.ID, &assistant.Name, &assistant.Model, &assistant.Type, &systemPrompts, &metadata)
        if err != nil {
            return types.AssistantList{}, err
        }

        assistant.SystemPrompts = systemPrompts

        // Unmarshal metadata JSON string into assistant.Metadata
        if metadata != "" {
            var meta types.Metadata
            err = json.Unmarshal([]byte(metadata), &meta)
            if err != nil {
                return types.AssistantList{}, fmt.Errorf("failed to unmarshal metadata: %v", err)
            }
            assistant.Metadata = &meta
        }

        assistants = append(assistants, assistant)
    }

    if err := rows.Err(); err != nil {
        return types.AssistantList{}, err
    }

    return types.AssistantList{Assistants: assistants}, nil
}