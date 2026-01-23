//go:build integration

package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/database/postgres"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

func TestMigrationsApplied(t *testing.T) {
	dbConn, cleanup := setupIntegrationDB(t)
	defer cleanup()

	assertTableExists(t, dbConn, "users")
	assertTableExists(t, dbConn, "medicine_categories")
	assertTableExists(t, dbConn, "notification_events")
	assertTableExists(t, dbConn, "support_chat_requests")
}

func TestUserAndProfileRepositories(t *testing.T) {
	dbConn, cleanup := setupIntegrationDB(t)
	defer cleanup()

	userRepo := NewUserRepository(dbConn)
	profileRepo := NewProfileRepository(dbConn)

	user := &db.User{Username: "0800000000", PasswordHash: "hash", Role: constants.RolePatient, IsActive: true, IsVerified: true}
	if err := userRepo.Create(context.Background(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	profile := &db.UserProfile{UserID: user.ID, FirstName: "A", LastName: "B", DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)}
	if err := profileRepo.Upsert(context.Background(), profile); err != nil {
		t.Fatalf("upsert profile: %v", err)
	}

	found, err := profileRepo.FindByUserID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("find profile: %v", err)
	}
	if found.FirstName != "A" {
		t.Fatalf("unexpected profile data")
	}
}

func TestNotificationRepository(t *testing.T) {
	dbConn, cleanup := setupIntegrationDB(t)
	defer cleanup()

	repo := NewNotificationRepository(dbConn)

	tpl := &db.NotificationTemplate{Code: "TEST_CODE", Title: "Title", Body: "Body", IsActive: true}
	if err := dbConn.Create(tpl).Error; err != nil {
		t.Fatalf("create template: %v", err)
	}

	userID := uuid.New()
	if err := dbConn.Create(&db.User{ID: userID, Username: "0810000000", PasswordHash: "hash", Role: constants.RolePatient, IsActive: true, IsVerified: true}).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	event := db.NotificationEvent{UserID: userID, TemplateCode: tpl.Code, ScheduledAt: time.Now().Add(-time.Minute), Status: constants.NotificationPending}
	if err := repo.CreateEvents(context.Background(), []db.NotificationEvent{event}); err != nil {
		t.Fatalf("create events: %v", err)
	}

	due, err := repo.ListDueForUpdate(context.Background(), time.Now().UTC(), 10)
	if err != nil {
		t.Fatalf("list due: %v", err)
	}
	if len(due) == 0 {
		t.Fatalf("expected due events")
	}
}

func setupIntegrationDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()
	cfg := loadDBConfigFromEnv()
	if cfg.Host == "" {
		t.Skip("DB_HOST not set")
	}

	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("find repo root: %v", err)
	}

	adminCfg := cfg
	adminCfg.Name = "postgres"
	adminDB, err := sql.Open("postgres", adminCfg.DSN())
	if err != nil {
		t.Fatalf("admin db connect: %v", err)
	}
	defer adminDB.Close()

	name := fmt.Sprintf("%s_test_%d", cfg.Name, time.Now().UnixNano())
	if _, err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(name))); err != nil {
		t.Skipf("create test database failed: %v", err)
	}

	cfg.Name = name
	m, err := migrate.New("file://"+filepath.Join(root, "migrations"), cfg.URL())
	if err != nil {
		t.Fatalf("migrate init: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("migrate up: %v", err)
	}

	dbConn, err := postgres.New(cfg)
	if err != nil {
		t.Fatalf("db open: %v", err)
	}

	cleanup := func() {
		sqlDB, _ := dbConn.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		_, _ = adminDB.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1", name)
		_, _ = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", pq.QuoteIdentifier(name)))
	}
	return dbConn, cleanup
}

func assertTableExists(t *testing.T, dbConn *gorm.DB, table string) {
	t.Helper()
	var exists *string
	row := dbConn.Raw("SELECT to_regclass(?)", "public."+table).Row()
	if err := row.Scan(&exists); err != nil || exists == nil {
		t.Fatalf("expected table %s", table)
	}
}

func loadDBConfigFromEnv() config.DBConfig {
	port := getEnvInt("DB_PORT", 5432)
	return config.DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		Name:     getEnv("DB_NAME", "stin_smart_care"),
		User:     getEnv("DB_USER", "stin"),
		Password: getEnv("DB_PASSWORD", "stin_pass"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func findRepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	current := wd
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(current, "migrations")); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("repo root not found")
}
