package handlers

import(
	"encoding/json"
	"fmt"
	"io"
	"github.com/go-playground/validator/v10"
    // "rag-demo/pkg/db"
    "context"
    "rag-demo/pkg/auth"
    "rag-demo/types"
	"net/http"
	"strings"
	"errors"
)

func decodeAndValidateJSON(body io.Reader, v interface{}) error {
    // Decode JSON
    if err := json.NewDecoder(body).Decode(v); err != nil {
        return fmt.Errorf("invalid JSON: %v", err)
    }

    // Validate struct
    validate := validator.New()
    if err := validate.Struct(v); err != nil {
        return fmt.Errorf("validation error: %v", err)
    }

    return nil
}

// ExtractAccessToken extracts the access token from the request.
// It checks the "access-token" header and cookie, and supports the "Bearer" prefix.
func ExtractAccessToken(r *http.Request) (string, error) {
    // Try to get the token from the "access-token" header
    authHeader := r.Header.Get("access-token")
    if authHeader != "" {
        // Check if the header value starts with "Bearer "
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
            return parts[1], nil
        }
        return authHeader, nil
    }

    // If not found in header, try to get the token from the cookies
    cookie, err := r.Cookie("access-token")
    if err == nil {
        return cookie.Value, nil
    }

    // Token not found in both header and cookies
    return "", errors.New("no access-token provided")
}

// Remove the creation of AuthService here
func GetSession(authService *auth.AuthServiceImpl, accessToken string) (types.Session, error) {
    // Retrieve the session from the access token
    session, err := authService.ValidateJWT(context.Background(), accessToken)
    if err != nil {
        return types.Session{}, err
    }

    return session, nil
}
