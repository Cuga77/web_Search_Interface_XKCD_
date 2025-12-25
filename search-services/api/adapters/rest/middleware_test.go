package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"yadro.com/course/api/adapters/auth"
	"yadro.com/course/api/adapters/rest"
)

type MockAuthMiddlewareHelper struct {
	mock.Mock
}

func (m *MockAuthMiddlewareHelper) Login(name, password string) (string, error) {
	args := m.Called(name, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthMiddlewareHelper) Validate(token string) (*auth.CustomClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.CustomClaims), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	mockAuth := new(MockAuthMiddlewareHelper)
	mw := rest.NewMiddleware(log, mockAuth)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := mw.AuthMiddleware(nextHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	mockAuth.On("Validate", "bad").Return(nil, errors.New("invalid")).Once()
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer bad")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	mockAuth.On("Validate", "good").Return(&auth.CustomClaims{User: "user"}, nil).Once()
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer good")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestConcurrencyMiddleware(t *testing.T) {
	mw := rest.NewMiddleware(log, nil)
	limit := 1

	blocker := make(chan struct{})
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-blocker
		w.WriteHeader(http.StatusOK)
	})

	handler := mw.ConcurrencyLimitMiddleware(limit, nextHandler)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}()

	time.Sleep(10 * time.Millisecond)

	go func() {
		defer wg.Done()
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	}()

	time.Sleep(10 * time.Millisecond)
	close(blocker)
	wg.Wait()
}

func TestRateLimitMiddleware(t *testing.T) {
	mw := rest.NewMiddleware(log, nil)
	rps := 1

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := mw.RateLimitMiddleware(rps, nextHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	req, _ = http.NewRequest("GET", "/", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
