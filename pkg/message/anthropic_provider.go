// anthropic_provider.go

package message

import (
    // "context"
    "encoding/json"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
    bedrockTypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
    "rag-demo/types"
)

type AnthropicProvider struct{}

func (p AnthropicProvider) BuildRequest(message types.MessageRequest) ([]byte, error) {
    msg := types.AnthropicMessage{
        Role:    "user",
        Content: message.Text,
    }
    requestData := types.AnthropicMessageRequest{
		// should probably not hardcode this
        AnthropicVersion: "bedrock-2023-05-31",
        MaxTokens:        3000,
		// will be retrieved from assistant table later
        System:           "You are a helpful assistant.",
        Messages:         []types.AnthropicMessage{msg},
    }
    return json.Marshal(requestData)
}

func (p AnthropicProvider) ProcessResponse(output *bedrockruntime.InvokeModelWithResponseStreamOutput) (string, error) {
    var combinedResult string
    for event := range output.GetStream().Events() {
        switch v := event.(type) {
        case *bedrockTypes.ResponseStreamMemberChunk:
            var data map[string]interface{}
            err := json.Unmarshal(v.Value.Bytes, &data)
            if err != nil {
                fmt.Println("Error unmarshaling chunk:", err)
                continue
            }
            if dataType, ok := data["type"].(string); ok {
                if dataType == "content_block_delta" {
                    if delta, ok := data["delta"].(map[string]interface{}); ok {
                        if text, ok := delta["text"].(string); ok {
                            fmt.Print(text)
                            combinedResult += text
                        }
                    }
                }
            }
        case *bedrockTypes.UnknownUnionMember:
            fmt.Println("unknown tag:", v.Tag)
        default:
            fmt.Println("union is nil or unknown type")
        }
    }
    return combinedResult, nil
}
