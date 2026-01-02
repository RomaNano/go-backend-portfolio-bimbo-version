package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/romanrey/go-backend-portfolio-bimbo-version/users-api-http/internal/handler"
	"github.com/romanrey/go-backend-portfolio-bimbo-version/users-api-http/internal/middleware"
	"github.com/romanrey/go-backend-portfolio-bimbo-version/users-api-http/internal/repository/postgres"
	"github.com/romanrey/go-backend-portfolio-bimbo-version/users-api-http/internal/service"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	
	userRepo := postgres.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.Create)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	rootHandler := middleware.RequestID(
		middleware.Logging(mux),
	)

	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", rootHandler))
}
