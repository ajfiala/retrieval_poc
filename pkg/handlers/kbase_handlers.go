package handlers

import (
	"fmt"
	"encoding/json"
	"net/http"
	"rag-demo/pkg/kbase"
	"rag-demo/types"
	"sync"
	"github.com/google/uuid"
	// "strings"
	"github.com/go-chi/chi/v5"
	// "github.com/go-playground/validator/v10"
	// "github.com/google/uuid"
)

func HandleCreateKbase(kbaseService kbase.KbaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newKbaseReq types.NewKbaseRequest
		err := decodeAndValidateJSON(r.Body, &newKbaseReq)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		resultCh := make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
		wg := &sync.WaitGroup{}
		wg.Add(1)

		kb := types.Kbase{
			ID: uuid.New(),
			Name: newKbaseReq.Name,
			Description: newKbaseReq.Description,
		}

		// Call CreateKbase with resultCh and wg
		go kbaseService.CreateKbase(r.Context(), kb, resultCh, wg)

		wg.Wait()             // Wait for the goroutine to finish
		result := <-resultCh  // Read the result from the channel


		if result.Success {
			kbase, ok := result.Data.(types.Kbase)
			if !ok {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kbase)
		} else {
			http.Error(w, "error creating kbase", http.StatusInternalServerError)
		}
	}
}

func HandleListKbases(kbaseService kbase.KbaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resultCh := make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go kbaseService.ListKbases(r.Context(), resultCh, wg)

		wg.Wait()             // Wait for the goroutine to finish
		result := <-resultCh  // Read the result from the channel

		fmt.Println("Result: ", result)

		if result.Success {
			kbases, ok := result.Data.(types.KbaseList)
			if !ok {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			fmt.Println("Kbases: ", kbases)
			json.NewEncoder(w).Encode(kbases)
		} else {
			fmt.Println("Error: ", result.Error)
			http.Error(w, "error listing kbases", http.StatusInternalServerError)
		}
	}
}

func HandleDeleteKbase(kbaseService kbase.KbaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kbaseID := chi.URLParam(r, "id")
		if kbaseID == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		kbID, err := uuid.Parse(kbaseID)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		resultCh := make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go kbaseService.DeleteKbase(r.Context(), kbID, resultCh, wg)

		wg.Wait()             // Wait for the goroutine to finish
		result := <-resultCh  // Read the result from the channel

		fmt.Println("HandleDeleteKbase result: ", result)

		
        if result.Success {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Kbase deleted successfully"))
        } else if result.Error == nil {
            // This handles the case where the kbase doesn't exist
            http.Error(w, "kbase not found", http.StatusNotFound)
        } else {
            http.Error(w, "error deleting kbase", http.StatusInternalServerError)
        }
	}
}