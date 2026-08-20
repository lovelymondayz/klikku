package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"klikku/internal/config"
	"klikku/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	return pool, nil
}

func RunMigrations(pool *pgxpool.Pool, dbName string) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	// Create migrations tracking table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Read migration files
	migrationsDir := "./migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".sql" {
			continue
		}
		if !strings.Contains(f.Name(), ".up.") {
			continue
		}

		// Extract version from filename (e.g., 000001_create_merchants.up.sql -> 1)
		var version int
		fmt.Sscanf(f.Name(), "%d", &version)

		// Check if already applied
		var applied bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&applied)
		if err != nil {
			return fmt.Errorf("check migration %d: %w", version, err)
		}
		if applied {
			continue
		}

		// Read and execute migration
		path := filepath.Join(migrationsDir, f.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", f.Name(), err)
		}

		log.Printf("Running migration: %s", f.Name())

		// Execute migration in a transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", f.Name(), err)
		}

		// Split by statements and execute each
		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "--") {
				continue
			}
			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return fmt.Errorf("execute %s: %w", f.Name(), err)
			}
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %d: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit %s: %w", f.Name(), err)
		}
	}

	log.Println("✅ Migrations complete")
	return nil
}

func SeedSuperAdmin(pool *pgxpool.Pool, cfg *config.Config) error {
	ctx := context.Background()

	var exists bool
	err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", cfg.SuperAdminEmail).Scan(&exists)
	if err != nil {
		return err
	}
	// Create default merchant if none exists
	var merchantID string
	err = pool.QueryRow(ctx, "SELECT id FROM merchants LIMIT 1").Scan(&merchantID)
	if err != nil {
		// No merchants exist, create one
		err = pool.QueryRow(ctx, "INSERT INTO merchants (business_name, slug, subscription_status) VALUES ($1, $2, $3) RETURNING id",
			"Default Merchant", "default", "active").Scan(&merchantID)
		if err != nil {
			return fmt.Errorf("create default merchant: %w", err)
		}
		log.Println("✅ Default merchant created")
	}

	if exists {
		log.Println("ℹ️  Super admin already exists")
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.SuperAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	_, err = pool.Exec(ctx, `INSERT INTO users (name, email, password_hash, role, merchant_id) VALUES ($1, $2, $3, $4, $5)`,
		"Super Admin", cfg.SuperAdminEmail, string(hash), models.RoleSuperAdmin, merchantID)
	if err != nil {
		return fmt.Errorf("create super admin: %w", err)
	}

	log.Println("✅ Super admin seeded")
	return nil
}
