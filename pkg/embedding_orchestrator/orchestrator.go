package orchestrator

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/google/uuid"

    "github.com/pgvector/pgvector-go"
    // "github.com/lib/pq"
    "rag-demo/pkg/index"
    "rag-demo/types"
)

type Orchestrator struct {
    bedrockService *index.BedrockRuntimeService
    dbService      types.KbaseEmbeddingsTableGateway
}

func NewOrchestrator(bedrockService *index.BedrockRuntimeService, dbService types.KbaseEmbeddingsTableGateway) *Orchestrator {
    return &Orchestrator{
        bedrockService: bedrockService,
        dbService:      dbService,
    }
}

func (o *Orchestrator) ProcessAndStoreEmbeddings(ctx context.Context, docText types.DocumentText, kbaseID uuid.UUID) error {
    for i, chunk := range docText.Chunks {
        // Prepare input for GetEmbeddings
        chunkDoc := types.DocumentText{
            Name:   docText.Name,
            Chunks: []string{chunk},
        }

        // Get embedding for the chunk
        embeddingOutput, err := o.bedrockService.GetEmbeddings(ctx, chunkDoc)
        if err != nil {
            return fmt.Errorf("error getting embeddings: %w", err)
        }

        // Parse the embedding output
        var embeddingResponse struct {
            Embedding []float32 `json:"embedding"`
        }
        
        err = json.Unmarshal(embeddingOutput.Body, &embeddingResponse)
        if err != nil {
            return fmt.Errorf("error unmarshaling embedding response: %w", err)
        }
        
        embeddingVec := pgvector.NewVector(embeddingResponse.Embedding)
        
        // Create the embedding record
        embeddingRecord := types.KbaseEmbedding{
            UUID:      uuid.New(),
            KbaseID:   kbaseID,
            ChunkID:   i,
            Content:   chunk,
            Embedding: embeddingVec,
            Metadata:  map[string]interface{}{"source": docText.Name},
        }

        // Store the embedding
        success, err := o.dbService.CreateEmbedding(ctx, embeddingRecord)
        if err != nil {
            return fmt.Errorf("error storing embedding: %w", err)
        }
        if !success {
            return fmt.Errorf("failed to store embedding for chunk %d", i)
        }
    }

    return nil
}