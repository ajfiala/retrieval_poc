package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"context"
	"fmt"
	"rag-demo/types"
	"github.com/google/uuid"
	"errors"
)


func (as *AuthServiceImpl) GenerateJWT(ctx context.Context, session types.Session) (string, error) {
    claims := map[string]interface{}{
        "userID": session.UserID.String(),
        "sessionID":   session.ID.String(),
    }
    _, tokenString, err := as.tokenAuth.Encode(claims)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}
func (as *AuthServiceImpl) ValidateJWT(ctx context.Context, tokenString string) (types.Session, error) {
	token, err := jwtauth.VerifyToken(as.tokenAuth, tokenString)
	if err != nil {
		fmt.Println("Error verifying token:", err)
		return types.Session{}, err
	}

	if token == nil {
		fmt.Println("Error: invalid token")
		return types.Session{}, errors.New("invalid token")
	}

	// Extract claims from the token
	claims, err := token.AsMap(ctx)
	if err != nil {
		fmt.Println("Error extracting claims from token:", err)
		return types.Session{}, err
	}

	userIDStr, ok := claims["userID"].(string)
	if !ok {
		fmt.Println("Error: invalid user ID in token")
		return types.Session{}, errors.New("invalid user ID in token")
	}

	sessionIDStr, ok := claims["sessionID"].(string)
	if !ok {
		fmt.Println("Error: invalid session ID in token")
		return types.Session{}, errors.New("invalid session ID in token")
	}

	// Parse the userID to a UUID
	_, err = uuid.Parse(userIDStr)
	if err != nil {
		fmt.Println("Error parsing user ID:", err)
		return types.Session{}, errors.New("invalid user ID format")
	}

	parsedSessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		fmt.Println("Error parsing session ID:", err)
		return types.Session{}, errors.New("invalid session ID format")
	}

	// Retrieve the user from the database using the existing UserGateway
	session, err := as.SessionGateway.GetSession(ctx, parsedSessionID)
	if err != nil {
		fmt.Printf("failed to retrieve sessionID: %s ", parsedSessionID)
		fmt.Println("with userID: ", userIDStr)
		fmt.Println("Error retrieving user from database:", err)
		return types.Session{}, errors.New("session not found in database")
	}

	// Create a new session object
	// session := types.Session{
	// 	ID:     parsedSessionID,
	// 	UserID: parsedUserID,
	// }

	return session, nil
}
