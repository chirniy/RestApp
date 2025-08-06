package data

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"serversTest2/internal/config"
	"serversTest2/internal/domain"
	"serversTest2/internal/migration"
)

func InitPostgresDB(cfg *config.Config) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", cfg.DataBaseURL)
	if err != nil {
		log.Fatal("DB open error:", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
		return nil, err
	}

	if err := migration.RunMigrations(db, "migrations"); err != nil {
		return nil, fmt.Errorf("migrations error: %w", err)
	}

	return db, nil
}

func InitInMemoryDB() (db map[uuid.UUID]domain.User, err error) {
	db = make(map[uuid.UUID]domain.User)
	return db, nil
}
