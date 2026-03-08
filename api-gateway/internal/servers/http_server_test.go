package servers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupServer() *HttpServer {
	gin.SetMode(gin.TestMode)

	wg := &sync.WaitGroup{}

	server := NewHttpServer(
		wg,
		nil, nil, nil, nil,
		nil, // userClient
		nil, // sessionClient
		nil, // messageClient
		// nil,
	)

	return server
}

func TestLogin_InvalidJSON(t *testing.T) {
	server := setupServer()

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_InvalidJSON(t *testing.T) {
	server := setupServer()

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMe_Unauthorized(t *testing.T) {
	server := setupServer()

	req := httptest.NewRequest(http.MethodGet, "/me", nil)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogout(t *testing.T) {
	server := setupServer()

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidate_InvalidJSON(t *testing.T) {
	server := setupServer()

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}