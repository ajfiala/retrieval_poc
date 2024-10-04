package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"context"
	"fmt"
	
	"rag-demo/types"
	"github.com/google/uuid"
	"errors"
)


func (as *AuthServiceImpl) GenerateJWT(ctx context.Context, user types.User) (string, error) {
    claims := map[string]interface{}{
        "userID": user.UserID.String(),
        "name":   user.Name,
    }
    _, tokenString, err := as.tokenAuth.Encode(claims)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}
func (as *AuthServiceImpl) ValidateJWT(ctx context.Context, tokenString string) (types.User, error) {
	token, err := jwtauth.VerifyToken(as.tokenAuth, tokenString)
	if err != nil {
		fmt.Println("Error verifying token:", err)
		return types.User{}, err
	}

	if token == nil {
		fmt.Println("Error: invalid token")
		return types.User{}, errors.New("invalid token")
	}

	// Extract claims from the token
	claims, err := token.AsMap(ctx)
	if err != nil {
		fmt.Println("Error extracting claims from token:", err)
		return types.User{}, err
	}

	userIDStr, ok := claims["userID"].(string)
	if !ok {
		fmt.Println("Error: invalid user ID in token")
		return types.User{}, errors.New("invalid user ID in token")
	}

	// Parse the userID to a UUID
	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		fmt.Println("Error parsing user ID:", err)
		return types.User{}, errors.New("invalid user ID format")
	}

	// Retrieve the user from the database using the existing UserGateway
	user, err := as.UserGateway.GetUser(ctx, parsedUserID)
	if err != nil {
		fmt.Println("userID: ", parsedUserID)
		fmt.Println("Error retrieving user from database:", err)
		return types.User{}, errors.New("user not found in database")
	}

	return user, nil
}
