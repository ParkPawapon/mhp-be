package handlers

import (
	"net/http"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHealthHandlers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("sqlite open failed: %v", err)
	}

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis failed: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	router := newTestRouter()
	handler := NewHealthHandler(db, rdb)
	router.GET("/healthz", handler.Healthz)
	router.GET("/readyz", handler.Readyz)

	resp := performRequest(router, http.MethodGet, "/healthz", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("healthz expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/readyz", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("readyz expected 200, got %d", resp.Code)
	}
}
