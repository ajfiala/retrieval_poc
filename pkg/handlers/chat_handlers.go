package handlers

import (
	// "fmt"
	// "encoding/json"
	// "net/http"
	// "rag-demo/pkg/message"
	// "rag-demo/types"
	// "sync"
	// "github.com/google/uuid"
	// "strings"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-playground/validator/v10"
	// "github.com/google/uuid"
)

// func HandleSendMessage(messageService message.MessageService) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var newMessageReq types.MessageRequest



// 		err := decodeAndValidateJSON(r.Body, &newMessageReq)
// 		if err != nil {
// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
// 			return
// 		}
		
// 		resultCh := make(types.ResultChannel, 1) 
// 		wg := &sync.WaitGroup{}
// 		wg.Add(1)

	
// 		go kbaseService.CreateKbase(r.Context(), kb, resultCh, wg)

// 		wg.Wait()             
// 		result := <-resultCh 


// 		if result.Success {
// 			kbase, ok := result.Data.(types.Kbase)
// 			if !ok {
// 				http.Error(w, "Internal server error", http.StatusInternalServerError)
// 				return
// 			}

// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(kbase)
// 		} else {
// 			http.Error(w, "error creating kbase", http.StatusInternalServerError)
// 		}
// 	}
// }
