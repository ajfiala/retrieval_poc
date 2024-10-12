
package message

import (
    "fmt"
)

func GetProvider(providerName string) (Provider, error) {
    switch providerName {
    case "anthropic":
        return AnthropicProvider{}, nil
    case "a121":
        return A121Provider{}, nil
    // Add other providers here
    default:
        return nil, fmt.Errorf("unknown provider: %s", providerName)
    }
}