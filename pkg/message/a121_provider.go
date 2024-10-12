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

type A121Provider struct{}

func (p A121Provider) BuildRequest(message types.MessageRequest) ([]byte, error) {
    msg := types.A121Message{
        Role:    "user",
        Content: message.Text,
    }
    requestData := types.A121MessageRequest{
        Messages:         []types.A121Message{msg},
        Number: 1,
    }
    return json.Marshal(requestData)
}

func (p A121Provider) ProcessResponse(output *bedrockruntime.InvokeModelWithResponseStreamOutput) error {
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

            if choices, ok := data["choices"].([]interface{}); ok {
                for _, choice := range choices {
                    if choiceMap, ok := choice.(map[string]interface{}); ok {
                        if delta, ok := choiceMap["delta"].(map[string]interface{}); ok {
                            // Check if 'content' exists in 'delta'
                            if content, ok := delta["content"].(string); ok {
                                fmt.Print(content)
                                combinedResult += content
                            }
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
    // Optionally, print the combined result
    // fmt.Println("\nCombined result:", combinedResult)
    return nil
}
