package main

import (
	"context"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"rag-demo/pkg/auth"
	"rag-demo/pkg/message"
	"rag-demo/pkg/kbase"
	"os"
	"rag-demo/pkg/db"
	"rag-demo/pkg/handlers"
)

func main() {
	// Set up database connection
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	err = db.RegisterType()
	if err != nil {
		log.Fatalf("Error registering type: %v", err)
	}
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Create user gateway
	userGateway := db.NewUserTableGateway(dbPool)

	// Create auth service
	authService := auth.NewAuthService(userGateway)

	// Create session gateway
	sessionGateway := db.NewSessionTableGateway(dbPool)

	// Create session service
	sessionService := message.NewSessionService(sessionGateway)

	// create kbase service 
	kbaseService := kbase.NewKbaseService(db.NewKbaseTableGateway(dbPool))

	// Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// Use the new handler that takes authService as an argument
	r.Post("/api/v1/signup", handlers.HandleCreateUser(authService))
	r.Post("/api/v1/validate", handlers.HandleValidateUser(authService))
	r.Get("/api/v1/user/{userID}", handlers.HandleGetUser(authService))
	r.Post("/api/v1/session", handlers.HandleCreateSession(sessionService))
	r.Post("/api/v1/kbase", handlers.HandleCreateKbase(kbaseService))
	r.Get("/api/v1/kbase", handlers.HandleListKbases(kbaseService))
	r.Delete("/api/v1/kbase/{id}", handlers.HandleDeleteKbase(kbaseService))

	// Start the server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}