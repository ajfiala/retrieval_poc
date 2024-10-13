package message

import (
    "context"
    "rag-demo/types"
    "github.com/google/uuid"
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

func (ms *MessageService) SendMessage(ctx context.Context, req types.MessageRequest) (*types.Message, error) {
    aiResponse, err := ms.BedrockService.InvokeModel(ctx, req)
    if err != nil {
        return nil, err
    }

    message := types.Message{
        ID:          uuid.New(),
        UserMessage: req.Text,
        AiMessage:   aiResponse,
    }

    _, err = ms.MessageGateway.StoreMessage(ctx, message)
    if err != nil {
        return nil, err
    }

    return &message, nil
}

func (ms *MessageService) RetrieveMessages(ctx context.Context, sessionID uuid.UUID) (types.MessageList, error) {
    return ms.MessageGateway.RetrieveMessages(ctx, sessionID)
}