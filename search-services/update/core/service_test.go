package core_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"yadro.com/course/update/core"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Add(ctx context.Context, comic core.Comics) error {
	args := m.Called(ctx, comic)
	return args.Error(0)
}

func (m *MockDB) Stats(ctx context.Context) (core.DBStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(core.DBStats), args.Error(1)
}

func (m *MockDB) Drop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDB) IDs(ctx context.Context) ([]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

type MockXKCD struct {
	mock.Mock
}

func (m *MockXKCD) Get(ctx context.Context, id int) (core.XKCDInfo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(core.XKCDInfo), args.Error(1)
}

func (m *MockXKCD) LastID(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
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

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) PublishUpdate() error {
	args := m.Called()
	return args.Error(0)
}

var log = slog.New(slog.NewTextHandler(os.Stderr, nil))

func TestUpdateStats(t *testing.T) {
	mockDB := new(MockDB)
	mockXKCD := new(MockXKCD)
	service, err := core.NewService(log, mockDB, mockXKCD, nil, nil, 1)
	assert.NoError(t, err)

	mockDB.On("Stats", mock.Anything).Return(core.DBStats{WordsTotal: 10}, nil).Once()
	mockXKCD.On("LastID", mock.Anything).Return(2, nil).Once()

	stats, err := service.Stats(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 10, stats.WordsTotal)
	assert.Equal(t, 2, stats.ComicsTotal)
}

func TestDrop(t *testing.T) {
	mockDB := new(MockDB)
	mockBus := new(MockEventBus)
	service, err := core.NewService(log, mockDB, nil, nil, mockBus, 1)
	assert.NoError(t, err)

	mockDB.On("Drop", mock.Anything).Return(nil).Once()
	mockBus.On("PublishUpdate").Return(nil).Once()

	err = service.Drop(context.Background())
	assert.NoError(t, err)
	mockBus.AssertExpectations(t)
}

func TestUpdate_AlreadyRunning(t *testing.T) {
	service, err := core.NewService(log, nil, nil, nil, nil, 1)
	assert.NoError(t, err)
	assert.Equal(t, core.StatusIdle, service.Status(context.Background()))
}

func TestUpdate(t *testing.T) {
	mockDB := new(MockDB)
	mockXKCD := new(MockXKCD)
	mockWords := new(MockWords)
	mockBus := new(MockEventBus)

	service, err := core.NewService(log, mockDB, mockXKCD, mockWords, mockBus, 1)
	assert.NoError(t, err)

	mockXKCD.On("LastID", mock.Anything).Return(2, nil).Once()

	mockDB.On("IDs", mock.Anything).Return([]int{1}, nil).Once()

	comic2 := core.XKCDInfo{
		ID: 2, URL: "url2", Title: "t", Alt: "a", Transcript: "tr", SafeTitle: "st",
	}
	mockXKCD.On("Get", mock.Anything, 2).Return(comic2, nil).Once()

	mockWords.On("Norm", mock.Anything, "a t tr").Return([]string{"kw"}, nil).Once()

	expectedComic := core.Comics{
		ID: 2, URL: "url2", Title: "t", Alt: "a", Transcript: "tr", SafeTitle: "st", Words: []string{"kw"},
	}
	mockDB.On("Add", mock.Anything, expectedComic).Return(nil).Once()

	mockBus.On("PublishUpdate").Return(nil).Once()

	err = service.Update(context.Background())
	assert.NoError(t, err)

	mockXKCD.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockWords.AssertExpectations(t)
	mockBus.AssertExpectations(t)
}
