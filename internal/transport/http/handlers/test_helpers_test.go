package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

const testRequestID = "req-test-1"

func newTestRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(constants.RequestIDKey, testRequestID)
		c.Next()
	})
	for _, mw := range middlewares {
		router.Use(mw)
	}
	return router
}

func withActor(role constants.Role, actorID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.SetActor(c, actorID, role)
		c.Next()
	}
}

func performRequest(router *gin.Engine, method, path string, payload any) *httptest.ResponseRecorder {
	var body *bytes.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	} else {
		body = bytes.NewReader(nil)
	}

	req, _ := http.NewRequest(method, path, body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

type envelopeMeta struct {
	Meta httpx.Meta `json:"meta"`
}

func strPtr(value string) *string {
	return &value
}
