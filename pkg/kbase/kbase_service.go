package kbase 

import (
	"context"
	"rag-demo/types"
	// google UUID package
	"github.com/google/uuid"
	"sync"
)

// KbaseService defines the interface for Kbase-related operations.
type KbaseService interface {
	CreateKbase(ctx context.Context,  kbase types.Kbase, resultCh types.ResultChannel, wg *sync.WaitGroup) 
	DeleteKbase(ctx context.Context, kbase_id uuid.UUID, resultCh types.ResultChannel, wg *sync.WaitGroup)
	ListKbases(ctx context.Context, resultCh types.ResultChannel, wg *sync.WaitGroup) 
}

type KbaseServiceImpl struct {
	KbaseGateway types.KbaseTableGateway
}

func NewKbaseService(KbaseGateway types.KbaseTableGateway) KbaseService {
	return &KbaseServiceImpl{KbaseGateway: KbaseGateway}
}

func (ks *KbaseServiceImpl) CreateKbase(ctx context.Context, kbase types.Kbase, resultCh types.ResultChannel, wg *sync.WaitGroup) {
    defer wg.Done()

    success, err := ks.KbaseGateway.CreateKbase(ctx, kbase)
    if err != nil || !success {
        resultCh <- types.Result{
            Data:    nil,
            Error:   err,
            Success: false,
        }
        return
    }

    resultCh <- types.Result{
        Data:    kbase,
        Error:   nil,
        Success: true,
    }
}

func (ks *KbaseServiceImpl) DeleteKbase(ctx context.Context, kbase_id uuid.UUID, resultCh types.ResultChannel, wg *sync.WaitGroup) {
	defer wg.Done()

	success, err := ks.KbaseGateway.DeleteKbase(ctx, kbase_id)
	if err != nil || !success {
		resultCh <- types.Result{
			Data:    nil,
			Error:   err,
			Success: false,
		}
		return
	}

	resultCh <- types.Result{
		Data:    nil,
		Error:   nil,
		Success: true,
	}
}

func (ks *KbaseServiceImpl) ListKbases(ctx context.Context, resultCh types.ResultChannel, wg *sync.WaitGroup) {
	defer wg.Done()

	kbases, err := ks.KbaseGateway.ListKbases(ctx)
	if err != nil {
		resultCh <- types.Result{
			Data:    nil,
			Error:   err,
			Success: false,
		}
		return
	}

	resultCh <- types.Result{
		Data:    kbases,
		Error:   nil,
		Success: true,
	}
}