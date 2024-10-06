package types 

import (
    "github.com/google/uuid"
	"github.com/lib/pq"
	"context"
)

// Kbase represents a knowledge base which can be used to provide context to an assistant for RAG.
type Kbase struct {
    ID            uuid.UUID         `json:"id"`
    Name          string            `json:"name"`     // Name of the knowledge base
    Description   string            `json:"description"`    // Model used by the assistant
}

type NewKbaseRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type KbaseEmbedding struct {
    ID        int                    `json:"id"`
    UUID      uuid.UUID              `json:"uuid"`
    KbaseID   uuid.UUID              `json:"kbase_id"`
    ChunkID   int                    `json:"chunk_id"`
    Content   string                 `json:"content"`
    Embedding pq.Float64Array        `json:"embedding" db:"embedding"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// KbaseList holds a collection of assistants.
type KbaseList struct {
    Kbases []Kbase `json:"kbases"` // List of assistants
}

type KbaseTableGateway interface {
	CreateKbase(ctx context.Context, kbase Kbase) (bool, error)
	GetKbase(ctx context.Context, kbaseId uuid.UUID) (Kbase, error)
	UpdateKbase(ctx context.Context, kbase Kbase) (bool, error)
    DeleteKbase(ctx context.Context, kbaseId uuid.UUID) (bool, error)
	ListKbases(ctx context.Context) (KbaseList, error)
}

type KbaseEmbeddingsTableGateway interface {
    CreateEmbedding(ctx context.Context, embedding KbaseEmbedding) (bool, error)
    // GetEmbedding(ctx context.Context, uuid uuid.UUID) (KbaseEmbedding, error)
}
