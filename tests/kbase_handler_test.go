package tests

import (
	"fmt"
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"net/http"
	"net/http/httptest"
	"rag-demo/pkg/db"
	"rag-demo/pkg/handlers"
	"rag-demo/pkg/kbase"
	"rag-demo/types"
	"testing"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateKbaseHandler(t *testing.T) {
    router := chi.NewRouter()
    testDBPool, err := pgxpool.New(context.Background(), "postgresql://myuser:mypassword@localhost:5432/goragdb")
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer testDBPool.Close()

    kbaseGateway := db.NewKbaseTableGateway(testDBPool)
    kbaseService := kbase.NewKbaseService(kbaseGateway)

    router.Post("/api/v1/kbase", handlers.HandleCreateKbase(kbaseService))


    // Prepare a new session request
    newKbaseRequest := types.NewKbaseRequest{
        Name: 	  "Test999",
		Description: "This is a test knowledge base",
    }

    // Encode the request as JSON
    body, err := json.Marshal(newKbaseRequest)
    if err != nil {
        t.Fatalf("Failed to marshal new kbase request: %v", err)
    }


    // Create a new HTTP request
    req, err := http.NewRequest("POST", "/api/v1/kbase", bytes.NewReader(body))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Record the response
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Fatalf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Decode the response body
	var kbase_response types.Kbase
    err = json.NewDecoder(rr.Body).Decode(&kbase_response)
    if err != nil {
        t.Fatalf("Failed to decode response body: %v", err)
    }

    // Verify the kbase response name
    if kbase_response.Name != newKbaseRequest.Name {
        t.Errorf("kbase has incorrect name: got %v, want %v", kbase_response.Name, newKbaseRequest.Name)
    }

	// verify kbase response ID is in db 
	k, err := kbaseGateway.GetKbase(context.Background(), kbase_response.ID)
	if k.ID != kbase_response.ID {
		t.Errorf("Kbase ID not found in db")
	}
	if err != nil {
		t.Errorf("Error fetching kbase from db")
	}

	// Verify the kbase response description
	if kbase_response.Description != newKbaseRequest.Description {
		t.Errorf("kbase has incorrect description: got %v, want %v", kbase_response.Description, newKbaseRequest.Description)
	}

	// Delete the kbase
	_, err = kbaseGateway.DeleteKbase(context.Background(), kbase_response.ID)

	if err != nil {
		t.Errorf("Failed to delete kbase: %v", err)
	}

}

func TestListKbaseHandler(t *testing.T) {
	router := chi.NewRouter()
	ctx := context.Background()
	testDBPool, err := pgxpool.New(context.Background(), "postgresql://myuser:mypassword@localhost:5432/goragdb")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer testDBPool.Close()

	kbaseGateway := db.NewKbaseTableGateway(testDBPool)
	kbaseService := kbase.NewKbaseService(kbaseGateway)

	// Test data
    testKbase := types.Kbase{
        ID:          uuid.New(),
        Name:        "Test Kbase12",
        Description: "Test description for Kbase",
    }
	testKbase2 := types.Kbase{
		ID:          uuid.New(),
		Name:        "Test Kbase13",
		Description: "Test description for Kbase",
	}

    // Create a new kbase so we can test list 
    resultCh := make(types.ResultChannel, 2) // Buffered channel to prevent deadlock
    wg := &sync.WaitGroup{}
    wg.Add(1)
    go kbaseService.CreateKbase(ctx, testKbase, resultCh, wg)
	wg.Wait()           
	result := <-resultCh  

	assert.True(t, result.Success, "Result should be successful")
	wg.Add(1)
	go kbaseService.CreateKbase(ctx, testKbase2, resultCh, wg)
    wg.Wait()             // Wait for the goroutine to finish
    result = <-resultCh  // Read the result from the channel
    assert.True(t, result.Success, "Result should be successful")


	router.Get("/api/v1/kbase", handlers.HandleListKbases(kbaseService))

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/api/v1/kbase", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response body
	var kbase_response types.KbaseList
	err = json.NewDecoder(rr.Body).Decode(&kbase_response)

	fmt.Println("Kbase Response: ", kbase_response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Verify the number of kbases
	if len(kbase_response.Kbases) < 1 {
		t.Errorf("Expected at least one kbase, got %v", len(kbase_response.Kbases))
	}

	// Delete the kbase
	_, err = kbaseGateway.DeleteKbase(context.Background(), testKbase.ID)

	if err != nil {
		t.Errorf("Failed to delete kbase: %v", err)
	}

	_, err = kbaseGateway.DeleteKbase(context.Background(), testKbase2.ID)

	if err != nil {
		t.Errorf("Failed to delete kbase: %v", err)
	}

}

func TestDeleteKbaseHandler(t *testing.T) {
	router := chi.NewRouter()
	ctx := context.Background()
	testDBPool, err := pgxpool.New(context.Background(), "postgresql://myuser:mypassword@localhost:5432/goragdb")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer testDBPool.Close()

	kbaseGateway := db.NewKbaseTableGateway(testDBPool)
	kbaseService := kbase.NewKbaseService(kbaseGateway)

	// Test data
	testKbase := types.Kbase{
		ID:          uuid.New(),
		Name:        "Test Kbase12",
		Description: "Test description for Kbase",
	}

	// Create a new kbase so we can test delete 
	resultCh := make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go kbaseService.CreateKbase(ctx, testKbase, resultCh, wg)
	wg.Wait()           
	result := <-resultCh  

	assert.True(t, result.Success, "Result should be successful")

	router.Delete("/api/v1/kbase/{id}", handlers.HandleDeleteKbase(kbaseService))

	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/kbase/%s", testKbase.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// // Decode the response body
	// var response types.Result
	// err = json.NewDecoder(rr.Body).Decode(&response)

	// fmt.Println("Response: ", response)
	// if err != nil {
	// 	t.Fatalf("Failed to decode response body: %v", err)
	// }

	// // Verify the response
	// if !response.Success {
	// 	t.Errorf("Expected success response, got %v", response.Error)
	// }

	// Delete the kbase
	_, err = kbaseGateway.DeleteKbase(context.Background(), testKbase.ID)

	if err != nil {
		t.Errorf("Failed to delete kbase: %v", err)
	}
}