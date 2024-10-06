package tests

import (
    "context"
    "testing"
    "sync"
    "rag-demo/pkg/db"
    "rag-demo/pkg/kbase"
    "rag-demo/types"
    "fmt"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/joho/godotenv"
    "github.com/stretchr/testify/assert"
)

// setupTestKbaseDB initializes a test database connection.
func setupTestKbaseDB(t *testing.T) *pgxpool.Pool {
    dbURL := "postgresql://myuser:mypassword@localhost:5432/goragdb"
    testDBPool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    return testDBPool
}

func TestKbaseService_CreateKbase(t *testing.T) {
    ctx := context.Background()
    dbPool := setupTestKbaseDB(t)
    defer dbPool.Close()

    // Load environment variables
    err := godotenv.Load("../.env")
    if err != nil {
        t.Log("No .env file found")
    }

    kbaseGateway := db.NewKbaseTableGateway(dbPool)
    kbaseService := kbase.NewKbaseService(kbaseGateway)

    // Test data
    testKbase := types.Kbase{
        ID:          uuid.New(),
        Name:        "Test Kbase12",
        Description: "Test description for Kbase",
    }

    // Create a new kbase
    resultCh := make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
    wg := &sync.WaitGroup{}
    wg.Add(1)

    go kbaseService.CreateKbase(ctx, testKbase, resultCh, wg)

    wg.Wait()             // Wait for the goroutine to finish
    result := <-resultCh  // Read the result from the channel

    fmt.Println("Result: ", result)
    assert.True(t, result.Success, "Result should be successful")
    assert.NoError(t, result.Error, "CreateKbase should not return an error")

    createdKbase, ok := result.Data.(types.Kbase)
    assert.True(t, ok, "Result Data should be of type types.Kbase")
    assert.NotNil(t, createdKbase, "Kbase should not be nil")
	assert.Equal(t, testKbase.ID, createdKbase.ID, "Kbase ID should match")
    assert.Equal(t, testKbase.Name, createdKbase.Name, "Kbase Name should match")
    assert.Equal(t, testKbase.Description, createdKbase.Description, "Kbase Description should match")

	 // Create a new kbase
	resultCh = make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
	wg = &sync.WaitGroup{}
	wg.Add(1)

	go kbaseService.DeleteKbase(ctx, createdKbase.ID, resultCh, wg)

	wg.Wait()             // Wait for the goroutine to finish
	result = <-resultCh  // Read the result from the channel

	fmt.Println("Result: ", result)
	assert.True(t, result.Success, "Result should be successful")
  
}

func TestKbaseService_ListKbases(t *testing.T) {
    ctx := context.Background()
    dbPool := setupTestKbaseDB(t)
    defer dbPool.Close()

    // Load environment variables
    err := godotenv.Load("../.env")
    if err != nil {
        t.Log("No .env file found")
    }

    kbaseGateway := db.NewKbaseTableGateway(dbPool)
    kbaseService := kbase.NewKbaseService(kbaseGateway)

    // Prepare test data
    testKbase1 := types.Kbase{
        ID:          uuid.New(),
        Name:        "Test123",
        Description: "Test description for Kbase 1",
    }

    testKbase2 := types.Kbase{
        ID:          uuid.New(),
        Name:        "Test456",
        Description: "Test description for Kbase 2",
    }

    // Create test kbases
    resultCh := make(types.ResultChannel, 2)
    wg := &sync.WaitGroup{}
    wg.Add(2)

    // Create first kbase
    go kbaseService.CreateKbase(ctx, testKbase1, resultCh, wg)

    // Create second kbase
    go kbaseService.CreateKbase(ctx, testKbase2, resultCh, wg)

    // Wait for both to finish
    wg.Wait()

    // Drain resultCh
    for i := 0; i < 2; i++ {
        result := <-resultCh
        assert.True(t, result.Success, "Result should be successful")
        assert.NoError(t, result.Error, "CreateKbase should not return an error")
    }

    // Now list kbases
    resultCh = make(types.ResultChannel, 1)
    wg = &sync.WaitGroup{}
    wg.Add(1)

    go kbaseService.ListKbases(ctx, resultCh, wg)

    wg.Wait()
    result := <-resultCh

    assert.True(t, result.Success, "Result should be successful")
    assert.NoError(t, result.Error, "ListKbases should not return an error")

    kbaseList, ok := result.Data.(types.KbaseList)
    assert.True(t, ok, "Result Data should be of type types.KbaseList")
    assert.NotNil(t, kbaseList, "KbaseList should not be nil")
    assert.GreaterOrEqual(t, len(kbaseList.Kbases), 2, "There should be at least 2 kbases")

    // Check that our test kbases are in the list
    var foundKbase1, foundKbase2 bool
    for _, kbase := range kbaseList.Kbases {
        if kbase.ID == testKbase1.ID {
            foundKbase1 = true
            assert.Equal(t, testKbase1.Name, kbase.Name, "Kbase1 Name should match")
            assert.Equal(t, testKbase1.Description, kbase.Description, "Kbase1 Description should match")
        }
        if kbase.ID == testKbase2.ID {
            foundKbase2 = true
            assert.Equal(t, testKbase2.Name, kbase.Name, "Kbase2 Name should match")
            assert.Equal(t, testKbase2.Description, kbase.Description, "Kbase2 Description should match")
        }
    }

    assert.True(t, foundKbase1, "Test Kbase 1 should be in the list")
    assert.True(t, foundKbase2, "Test Kbase 2 should be in the list")

	// Deletion of testKbase1
    resultCh = make(types.ResultChannel, 1)
    wg = &sync.WaitGroup{}
    wg.Add(1)

    go kbaseService.DeleteKbase(ctx, testKbase1.ID, resultCh, wg)

    wg.Wait() // Wait for the goroutine to finish
    result = <-resultCh
    assert.True(t, result.Success, "Deletion of Kbase1 should be successful")

    // Deletion of testKbase2
    wg.Add(1)

    go kbaseService.DeleteKbase(ctx, testKbase2.ID, resultCh, wg)

    wg.Wait() // Wait for the goroutine to finish
    result = <-resultCh
    assert.True(t, result.Success, "Deletion of Kbase2 should be successful")
}
