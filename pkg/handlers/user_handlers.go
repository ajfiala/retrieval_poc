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


// In handlers/user_handlers.go
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
            createUserResult, ok := result.Data.(types.CreateUserResult)
            if !ok {
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            w.Header().Set("access-token", "Bearer "+ createUserResult.Token)

            http.SetCookie(w, &http.Cookie{
                Name:     "access-token",
                Value:    createUserResult.Token,
                Path:     "/",
                HttpOnly: true,
                Secure:   true, 
                SameSite: http.SameSiteLaxMode,
            })

            json.NewEncoder(w).Encode(createUserResult.User)
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

		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}


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

        token, err := ExtractAccessToken(r) 
		if err != nil {
			http.Error(w, "Invalid access-token", http.StatusUnauthorized)
			return
		}

        user, err := authService.ValidateJWT(r.Context(), token)
        if err != nil {
            http.Error(w, "Invalid access-token", http.StatusUnauthorized)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}