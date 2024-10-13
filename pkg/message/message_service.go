package message

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"rag-demo/types"
	"sync"
)

type BedrockRuntimeServicer interface {
	InvokeModel(ctx context.Context, message types.MessageRequest) (string, error)
}

type MessageService struct {
	MessageGateway types.MessageTableGateway
	BedrockService BedrockRuntimeServicer
}

func NewMessageService(gateway types.MessageTableGateway, bedrock BedrockRuntimeServicer) *MessageService {
	return &MessageService{
		MessageGateway: gateway,
		BedrockService: bedrock,
	}
}

func (ms *MessageService) SendMessage(ctx context.Context,
    req types.MessageRequest,
    session types.Session,
    resultCh types.ResultChannel,
    wg *sync.WaitGroup) {

    defer wg.Done()

    fmt.Println("calling BedrockService.InvokeModel...")
    aiResponse, err := ms.BedrockService.InvokeModel(ctx, req)
    if err != nil {
        resultCh <- types.Result{
            Data:    nil,
            Error:   err,
            Success: false,
        }
        return 
    }

    message := types.Message{
        ID:          uuid.New(),
        UserMessage: req.Text,
        AiMessage:   aiResponse,
        SessionId:   session.ID,
        UserId:     session.UserID,
    }

    fmt.Println("message: ", message)

    fmt.Println("calling MessageGateway.StoreMessage...")
    ok, err := ms.MessageGateway.StoreMessage(ctx, message, session)
    if err != nil {
        resultCh <- types.Result{
            Data:    nil,
            Error:   err,
            Success: false,
        }
        return
    }
    if !ok {
        resultCh <- types.Result{
            Data:    nil,
            Error:   nil,
            Success: false,
        }
        return 
    }

    resultCh <- types.Result{
        Data:    message,
        Error:   nil,
        Success: true,
    }
}


func (ms *MessageService) RetrieveMessages(ctx context.Context, sessionID uuid.UUID) (types.MessageList, error) {
	return ms.MessageGateway.RetrieveMessages(ctx, sessionID)
}
