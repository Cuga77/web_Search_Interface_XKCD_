package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"yadro.com/course/api/adapters/auth"
	"yadro.com/course/api/adapters/rest"
	"yadro.com/course/api/core"
)

type MockPinger struct {
	mock.Mock
}

func (m *MockPinger) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockUpdater struct {
	mock.Mock
}

func (m *MockUpdater) Update(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUpdater) Stats(ctx context.Context) (core.UpdateStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(core.UpdateStats), args.Error(1)
}

func (m *MockUpdater) Status(ctx context.Context) (core.UpdateStatus, error) {
	args := m.Called(ctx)
	return args.Get(0).(core.UpdateStatus), args.Error(1)
}

func (m *MockUpdater) Drop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockSearcher struct {
	mock.Mock
}

func (m *MockSearcher) Search(ctx context.Context, phrase string, limit int) (core.SearchResult, error) {
	args := m.Called(ctx, phrase, limit)
	return args.Get(0).(core.SearchResult), args.Error(1)
}

func (m *MockSearcher) ISearch(ctx context.Context, phrase string, limit int) (core.SearchResult, error) {
	args := m.Called(ctx, phrase, limit)
	return args.Get(0).(core.SearchResult), args.Error(1)
}

type MockAuthorizer struct {
	mock.Mock
}

func (m *MockAuthorizer) Login(name, password string) (string, error) {
	args := m.Called(name, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthorizer) Validate(tokenString string) (*auth.CustomClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.CustomClaims), args.Error(1)
}

var log = slog.New(slog.NewTextHandler(os.Stderr, nil))

func TestPingHandler(t *testing.T) {
	mockPinger := new(MockPinger)
	mockPinger.On("Ping", mock.Anything).Return(nil).Once()

	pingers := map[string]core.Pinger{
		"test": mockPinger,
	}

	handler := rest.NewPingHandler(log, pingers)

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["replies"]["test"])

	mockPinger.AssertExpectations(t)
}

func TestUpdateHandler(t *testing.T) {
	mockUpdater := new(MockUpdater)
	mockUpdater.On("Update", mock.Anything).Return(nil).Once()

	handler := rest.NewUpdateHandler(log, mockUpdater)

	req, _ := http.NewRequest(http.MethodPost, "/update", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockUpdater.AssertExpectations(t)
}

func TestUpdateHandler_InProgress(t *testing.T) {
	mockUpdater := new(MockUpdater)
	mockUpdater.On("Update", mock.Anything).Return(errors.New("update in progress")).Once()

	handler := rest.NewUpdateHandler(log, mockUpdater)

	req, _ := http.NewRequest(http.MethodPost, "/update", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)
	mockUpdater.AssertExpectations(t)
}

func TestUpdateHandler_Error(t *testing.T) {
	mockUpdater := new(MockUpdater)
	mockUpdater.On("Update", mock.Anything).Return(errors.New("internal logic error")).Once()

	handler := rest.NewUpdateHandler(log, mockUpdater)

	req, _ := http.NewRequest(http.MethodPost, "/update", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockUpdater.AssertExpectations(t)
}

func TestSearchHandler(t *testing.T) {
	mockSearcher := new(MockSearcher)
	expectedResult := core.SearchResult{
		Total: 1,
		Comics: []core.Comic{
			{ID: 1, URL: "url"},
		},
	}
	mockSearcher.On("Search", mock.Anything, "test", 10).Return(expectedResult, nil).Once()

	handler := rest.NewSearchHandler(log, mockSearcher)

	req, _ := http.NewRequest(http.MethodGet, "/search?phrase=test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var result core.SearchResult
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)

	mockSearcher.AssertExpectations(t)
}

func TestSearchHandler_BagParams(t *testing.T) {
	mockSearcher := new(MockSearcher)
	handler := rest.NewSearchHandler(log, mockSearcher)

	req, _ := http.NewRequest(http.MethodGet, "/search", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	req, _ = http.NewRequest(http.MethodGet, "/search?phrase=test&limit=abc", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginHandler(t *testing.T) {
	mockAuth := new(MockAuthorizer)
	mockAuth.On("Login", "admin", "password").Return("token", nil).Once()

	handler := rest.NewLoginHandler(log, mockAuth)

	reqBody := bytes.NewBufferString(`{"name":"admin", "password":"password"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", reqBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "token", rr.Body.String())
	mockAuth.AssertExpectations(t)
}

func TestLoginHandler_Fail(t *testing.T) {
	mockAuth := new(MockAuthorizer)
	mockAuth.On("Login", "admin", "wrong").Return("", errors.New("auth failed")).Once()

	handler := rest.NewLoginHandler(log, mockAuth)

	reqBody := bytes.NewBufferString(`{"name":"admin", "password":"wrong"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", reqBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockAuth.AssertExpectations(t)
}
