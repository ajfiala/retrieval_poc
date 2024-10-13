package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rag-demo/pkg/auth"
	"rag-demo/pkg/message"
	// "rag-demo/pkg/session"
	"rag-demo/types"
	"sync"
	// "github.com/google/uuid"
	// "strings"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-playground/validator/v10"
	// "github.com/google/uuid"
)

func HandleSendMessage(authService *auth.AuthServiceImpl, messageService *message.MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newMessageReq types.MessageRequest

		fmt.Println("validating JSON...")

		token, err := ExtractAccessToken(r) 
		if err != nil {
			http.Error(w, "No access-token provided", http.StatusUnauthorized)
			return
		}

		session, err := GetSession(authService, token)
        if err != nil {
            fmt.Println("error getting session: ", err)
            http.Error(w, "Invalid access-token", http.StatusUnauthorized)
            return
        }
		

		err = decodeAndValidateJSON(r.Body, &newMessageReq)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		resultCh := make(types.ResultChannel, 1) 
		wg := &sync.WaitGroup{}
		wg.Add(1)

		fmt.Println("calling messageService.SendMessage...")
		go messageService.SendMessage(r.Context(), newMessageReq, session, resultCh, wg)
	

		wg.Wait()             
		result := <-resultCh 

		fmt.Println("result: ", result)
		if result.Success {
			msg, ok := result.Data.(types.Message)
			if !ok {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(msg)
		} else {
			http.Error(w, "error creating kbase", http.StatusInternalServerError)
		}
	}
}
