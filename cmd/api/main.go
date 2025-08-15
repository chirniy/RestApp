package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"serversTest2/internal/config"
	"serversTest2/internal/data"
	"serversTest2/internal/handlers"
	"serversTest2/internal/middleware"
	"serversTest2/internal/repository/postgres"
	"serversTest2/internal/usecase"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
		return
	}
	log.Println("Starting server at " + cfg.Port)
	db, err := data.InitPostgresDB(cfg)
	if err != nil {
		log.Fatal("DB table create error:", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("DB close error:", err)
			return
		}
	}(db)
	repo := postgres.NewPostgresUserRepo(db)
	uc := usecase.NewUserUsecase(repo)
	handler := handlers.NewUserHandler(uc)
	r := mux.NewRouter()
	//users = make(map[uuid.UUID]user)
	r.HandleFunc("/users", handler.HomeHandler).Methods("GET", "OPTIONS", "POST")
	r.Handle("/users/{id}", middleware.JWTMiddleware(http.HandlerFunc(handler.UsersHandler))).Methods("GET", "OPTIONS", "PUT", "DELETE", "PATCH")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST", "OPTIONS")
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		return
	}
}
