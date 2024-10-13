package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rag-demo/pkg/auth"
	// "rag-demo/pkg/message"
	"rag-demo/pkg/assistant"
	"rag-demo/types"
	"sync"
	"github.com/google/uuid"
	// "strings"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-playground/validator/v10"
	// "github.com/google/uuid"
)


func HandleCreateAssistant(authService *auth.AuthServiceImpl, assistantService assistant.AssistantService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newAssistant types.Assistant

		token, err := ExtractAccessToken(r)
		if err != nil {
			http.Error(w, "No access-token provided", http.StatusUnauthorized)
			return
		}

		// validate session 
		_, err = GetSession(authService, token)
		if err != nil {
			fmt.Println("error getting session: ", err)
			http.Error(w, "Invalid access-token", http.StatusUnauthorized)
			return
		}

		err = decodeAndValidateJSON(r.Body, &newAssistant)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newAssistant.ID = uuid.New()

		resultCh := make(types.ResultChannel, 1)
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go assistantService.CreateAssistant(r.Context(), newAssistant, resultCh, wg)

		result := <-resultCh
		wg.Wait()

		if result.Success {
			assistant, ok := result.Data.(types.Assistant)
			if !ok {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(assistant)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleListAssistants(authService *auth.AuthServiceImpl, assistantService assistant.AssistantService) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		token, err := ExtractAccessToken(r)
		if err != nil {
			http.Error(w, "No access-token provided", http.StatusUnauthorized)
			return
		}

		// validate session
		_, err = GetSession(authService, token)
		if err != nil {
			fmt.Println("error getting session: ", err)
			http.Error(w, "Invalid access-token", http.StatusUnauthorized)
			return
		}

		resultCh := make(types.ResultChannel, 1)
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go assistantService.ListAssistants(r.Context(), resultCh, wg)

		result := <-resultCh
		wg.Wait()

		if result.Success {
			assistants, ok := result.Data.(types.AssistantList)
			if !ok {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(assistants)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
	}
}