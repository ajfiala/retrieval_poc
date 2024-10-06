package db

import (
	"context"
	"rag-demo/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserTableGatewayImpl is the implementation of UserTableGateway using pgxpool.
type KbaseTableGatewayImpl struct {
	Pool *pgxpool.Pool
}

// NewUserTableGateway creates a new instance of UserTableGatewayImpl.
func NewKbaseTableGateway(pool *pgxpool.Pool) types.KbaseTableGateway {
	return &KbaseTableGatewayImpl{Pool: pool}
}

// CreateKbase creates a new knowledge base in the kbase table of the Postgres db
func (k *KbaseTableGatewayImpl) CreateKbase(ctx context.Context, kbase types.Kbase) (bool, error) {
	_, err := k.Pool.Exec(ctx, "INSERT INTO kbase (uuid, name, description) VALUES ($1, $2, $3)", kbase.ID, kbase.Name, kbase.Description)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateKbase updates an existing knowledge base in the kbase table of the Postgres db
func (k *KbaseTableGatewayImpl) UpdateKbase(ctx context.Context, kbase types.Kbase) (bool, error) {
	_, err := k.Pool.Exec(ctx, "UPDATE kbase SET name = $1, description = $2 WHERE uuid = $3", kbase.Name, kbase.Description, kbase.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetKbase retrieves a knowledge base from the kbase table of the Postgres db
func (k *KbaseTableGatewayImpl) GetKbase(ctx context.Context, kbaseId uuid.UUID) (types.Kbase, error) {
	var kbase types.Kbase
	err := k.Pool.QueryRow(ctx, "SELECT uuid, name, description FROM kbase WHERE uuid = $1", kbaseId).Scan(&kbase.ID, &kbase.Name, &kbase.Description)
	if err != nil {
		return types.Kbase{}, err
	}
	return kbase, nil
}

// ListKbases retrieves all knowledge bases from the kbase table of the Postgres db
func (k *KbaseTableGatewayImpl) ListKbases(ctx context.Context) (types.KbaseList, error) {
	rows, err := k.Pool.Query(ctx, "SELECT uuid, name, description FROM kbase")
	if err != nil {
		return types.KbaseList{}, err
	}
	defer rows.Close()

	var kbases []types.Kbase
	for rows.Next() {
		var kbase types.Kbase
		err := rows.Scan(&kbase.ID, &kbase.Name, &kbase.Description)
		if err != nil {
			return types.KbaseList{}, err
		}
		kbases = append(kbases, kbase)
	}

	return types.KbaseList{Kbases: kbases}, nil
}

// DeleteKbase deletes a knowledge base from the kbase table of the Postgres db and the associated embeddings
func (k *KbaseTableGatewayImpl) DeleteKbase(ctx context.Context, kbaseId uuid.UUID) (bool, error) {
	tx, err := k.Pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	// first, check if it exists 
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM kbase WHERE uuid = $1)", kbaseId).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	_, err = tx.Exec(ctx, "DELETE FROM kbase_embeddings WHERE kbase_id = $1", kbaseId)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(ctx, "DELETE FROM kbase WHERE uuid = $1", kbaseId)
	if err != nil {
		return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}