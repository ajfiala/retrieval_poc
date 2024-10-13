package message

import (
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"rag-demo/types"
)

type Provider interface {
	BuildRequest(message types.MessageRequest) ([]byte, error)
	ProcessResponse(output *bedrockruntime.InvokeModelWithResponseStreamOutput) (string, error)
}