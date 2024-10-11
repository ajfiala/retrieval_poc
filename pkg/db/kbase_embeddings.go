package db

import (
    "context"
    "encoding/json"
    // "time"

    "rag-demo/types"
    // "github.com/pgvector/pgvector-go"
    // pgxvector "github.com/pgvector/pgvector-go/pgx"
    // "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type KbaseEmbeddingsTableGatewayImpl struct {
    Pool *pgxpool.Pool
}

func NewKbaseEmbeddingsTableGateway(pool *pgxpool.Pool) types.KbaseEmbeddingsTableGateway {
    return &KbaseEmbeddingsTableGatewayImpl{Pool: pool}
}

func (k *KbaseEmbeddingsTableGatewayImpl) CreateEmbedding(ctx context.Context, embedding types.KbaseEmbedding) (bool, error) {
    // Marshal Metadata to JSON
    metadataJSON, err := json.Marshal(embedding.Metadata)
    if err != nil {
        return false, err
    }


    // Insert into the database
    _, err = k.Pool.Exec(ctx, `
        INSERT INTO kbase_embeddings (uuid, kbase_id, chunk_id, content, embedding, metadata)
        VALUES ($1, $2, $3, $4, $5, $6)
    `,
        embedding.UUID,
        embedding.KbaseID,
        embedding.ChunkID,
        embedding.Content,
        embedding.Embedding, // Now a pgvector.Vector
        metadataJSON,
    )
    if err != nil {
        return false, err
    }
    return true, nil
}

// Implement GetEmbedding and other methods as needed
