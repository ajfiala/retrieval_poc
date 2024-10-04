package handlers

import (
	"encoding/json"
	"net/http"
	"rag-demo/pkg/auth"
	"rag-demo/types"
	"sync"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)


func HandleCreateUser(authService auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser types.NewUserRequest
		err := decodeAndValidateJSON(r.Body, &newUser)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resultCh := make(types.ResultChannel)
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go authService.CreateUser(r.Context(), newUser.Name, resultCh, wg)

		result := <-resultCh
		wg.Wait()

		if result.Success {
			user, ok := result.Data.(types.User)
			if !ok {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			// make jwt token
			token, err := authService.GenerateJWT(r.Context(), user)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}


			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("access-token", "Bearer "+ token)

			// Set the token as a cookie
            http.SetCookie(w, &http.Cookie{
                Name:     "access-token",
                Value:    token,
                Path:     "/",
                HttpOnly: true,
                Secure:   true, // Set to true if your site uses HTTPS
                SameSite: http.SameSiteLaxMode,
            })

			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleGetUser(authService auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			http.Error(w, "Missing user_id", http.StatusBadRequest)
			return
		}

		// Validate the user_id
		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		// userID is already validated and parsed

		resultCh := make(types.ResultChannel)
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go authService.GetUser(r.Context(), parsedUserID, resultCh, wg)

		result := <-resultCh
		wg.Wait()

		if result.Success {
			user, ok := result.Data.(types.User)
			if !ok {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		} else {
			if result.Error != nil && result.Error.Error() == "user not found" {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func HandleValidateUser(authService auth.AuthService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var token string

        // First, try to get the token from the "access-token" header
        token, err := ExtractAccessToken(r) 
		if err != nil {
			http.Error(w, "Invalid access-token", http.StatusUnauthorized)
			return
		}

        // Validate the token
        user, err := authService.ValidateJWT(r.Context(), token)
        if err != nil {
            http.Error(w, "Invalid access-token", http.StatusUnauthorized)
            return
        }

        // Token is valid; return user info as JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}