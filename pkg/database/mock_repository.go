package database

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ExecuteQuery(ctx context.Context, hostname string, startTime, endTime time.Time) (time.Duration, error) {
	args := m.Called(ctx, hostname, startTime, endTime)
	return time.Duration(1 * time.Millisecond), args.Error(0)
}

func (m *MockRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRepository) Close() {
	m.Called()
}
