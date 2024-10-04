package handlers

import (
	"encoding/json"
	"net/http"
	"rag-demo/pkg/message"
	"rag-demo/types"
	"sync"
	// "strings"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-playground/validator/v10"
	// "github.com/google/uuid"
)


func HandleCreateSession(sessionService message.SessionService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var newSession types.NewSessionRequest
        err := decodeAndValidateJSON(r.Body, &newSession)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

        // Parse userID from string to uuid.UUID
        parsedUserID := newSession.UserID

        resultCh := make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
        wg := &sync.WaitGroup{}
        wg.Add(1)

        // Call CreateSession with resultCh and wg
        go sessionService.CreateSession(r.Context(), parsedUserID, resultCh, wg)

        wg.Wait()             // Wait for the goroutine to finish
        result := <-resultCh  // Read the result from the channel

        if result.Success {
            session, ok := result.Data.(types.Session)
            if !ok {
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(session)
        } else {
            http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        }
    }
}