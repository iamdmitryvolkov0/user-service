package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"user-srv/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) RunMigrations() error {
	goose.SetBaseFS(nil)
	if err := goose.Up(m.db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")
	return nil
}

func (m *Migrator) Seed() error {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count users: %v", err)
	}

	if count >= 5 {
		log.Println("Seeder skipped: enough users already exist")
		return nil
	}

	users := []struct {
		name     string
		email    string
		password string
	}{
		{"Alice", "alice@example.com", "pass123"},
		{"Bob", "bob@example.com", "pass456"},
		{"Charlie", "charlie@example.com", "pass789"},
		{"Dave", "dave@example.com", "pass101"},
		{"Eve", "eve@example.com", "pass202"},
	}

	for i := count; i < 5; i++ {
		hashedPassword, err := hashPassword(users[i].password)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %v", users[i].name, err)
		}
		query := "INSERT INTO users (name, email, password, created_at) VALUES ($1, $2, $3, $4)"
		_, err = m.db.ExecContext(context.Background(), query, users[i].name, users[i].email, hashedPassword, time.Now())
		if err != nil {
			return fmt.Errorf("failed to seed user %s: %v", users[i].name, err)
		}
	}

	log.Printf("Seeder completed: added %d users", 5-count)
	return nil
}

func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return db, nil
}

func InitDB() *sql.DB {
	cfg := config.LoadConfig()
	db, err := ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	migrator := NewMigrator(db)
	if err := migrator.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	if err := migrator.Seed(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	return db
}
