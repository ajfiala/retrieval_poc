package assistant

import(
	"context"
	"rag-demo/types"
	// "github.com/google/uuid"
	"sync"
)

type AssistantService interface {
	CreateAssistant(ctx context.Context,  assistant types.Assistant, resultCh types.ResultChannel, wg *sync.WaitGroup) 
	ListAssistants(ctx context.Context, resultCh types.ResultChannel, wg *sync.WaitGroup) 
}

type AssistantServiceImpl struct {
	AssistantGateway types.AssistantTableGateway
}

func NewAssistantService(AssistantGateway types.AssistantTableGateway) AssistantService {
	return &AssistantServiceImpl{AssistantGateway: AssistantGateway}
}

func (as *AssistantServiceImpl) CreateAssistant(ctx context.Context,
	 assistant types.Assistant,
	  resultCh types.ResultChannel,
	   wg *sync.WaitGroup){
	defer wg.Done()
	success, err := as.AssistantGateway.CreateAssistant(ctx, assistant)
	if err != nil || !success {
		resultCh <- types.Result{
			Data: types.Assistant{},
			Error: err,
			Success: false,
		}
		return
	}
	resultCh <- types.Result{
		Data: assistant,
		Error: nil,
		Success: true,
	}
}

func (as *AssistantServiceImpl) ListAssistants(ctx context.Context, resultCh types.ResultChannel, wg *sync.WaitGroup) {
	defer wg.Done()
	assistants, err := as.AssistantGateway.ListAssistants(ctx)
	if err != nil {
		resultCh <- types.Result{
			Data:    nil,
			Error:   err,
			Success: false,
		}
		return
	}
	resultCh <- types.Result{
		Data:    assistants,
		Error:   nil,
		Success: true,
	}
}