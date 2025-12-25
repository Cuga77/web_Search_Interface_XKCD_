package core_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"yadro.com/course/search/core"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Search(ctx context.Context, keywords []string, limit int) ([]core.Comic, int64, error) {
	args := m.Called(ctx, keywords, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]core.Comic), args.Get(1).(int64), args.Error(2)
}

func (m *MockDB) Scan(ctx context.Context) ([]core.Comic, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]core.Comic), args.Error(1)
}

type MockWords struct {
	mock.Mock
}

func (m *MockWords) Norm(ctx context.Context, phrase string) ([]string, error) {
	args := m.Called(ctx, phrase)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

var log = slog.New(slog.NewTextHandler(os.Stderr, nil))

func TestSearch(t *testing.T) {
	mockDB := new(MockDB)
	mockWords := new(MockWords)
	service := core.NewService(log, mockDB, mockWords)

	mockWords.On("Norm", mock.Anything, "fail").Return(nil, errors.New("norm error")).Once()
	_, err := service.Search(context.Background(), "fail", 10)
	assert.Error(t, err)

	mockWords.On("Norm", mock.Anything, "empty").Return([]string{}, nil).Once()
	res, err := service.Search(context.Background(), "empty", 10)
	assert.NoError(t, err)
	assert.Empty(t, res.Comics)

	mockWords.On("Norm", mock.Anything, "test").Return([]string{"test"}, nil).Once()
	mockDB.On("Search", mock.Anything, []string{"test"}, 10).Return([]core.Comic{{ID: 1}}, int64(1), nil).Once()
	res, err = service.Search(context.Background(), "test", 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Comics))
	assert.Equal(t, int64(1), res.Total)

	mockWords.On("Norm", mock.Anything, "dbfail").Return([]string{"dbfail"}, nil).Once()
	mockDB.On("Search", mock.Anything, []string{"dbfail"}, 10).Return(nil, int64(0), errors.New("db error")).Once()
	_, err = service.Search(context.Background(), "dbfail", 10)
	assert.Error(t, err)
}

func TestISearch(t *testing.T) {
	mockDB := new(MockDB)
	mockWords := new(MockWords)
	service := core.NewService(log, mockDB, mockWords)

	mockDB.On("Scan", mock.Anything).Return([]core.Comic{
		{ID: 1, Keywords: []string{"test"}},
		{ID: 2, Keywords: []string{"foo"}},
	}, nil).Once()
	err := service.BuildIndex(context.Background())
	assert.NoError(t, err)

	mockWords.On("Norm", mock.Anything, "test").Return([]string{"test"}, nil).Once()
	res, err := service.ISearch(context.Background(), "test", 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Comics))
	assert.Equal(t, int64(1), res.Comics[0].ID)

	mockWords.On("Norm", mock.Anything, "bar").Return([]string{"bar"}, nil).Once()
	res, err = service.ISearch(context.Background(), "bar", 10)
	assert.NoError(t, err)
	assert.Empty(t, res.Comics)

	mockWords.On("Norm", mock.Anything, " ").Return([]string{}, nil).Once()
	res, err = service.ISearch(context.Background(), " ", 10)
	assert.NoError(t, err)
	assert.Empty(t, res.Comics)
}
