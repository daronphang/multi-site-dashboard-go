package ucmock

import (
	"errors"
	repomock "multi-site-dashboard-go/internal/repository/mock"
	"sync"
	"testing"
)

type UseCaseMocks struct {
	Repository * repomock.MockExtQuerier
	EventPublisher *MockEventPublisher
	Broadcaster *MockBroadcaster
}

var (
	ucm *UseCaseMocks
	syncOnce sync.Once
	ErrDB = errors.New("some error from db")
	ErrEventPublisher = errors.New("some error from publishing event")
	ErrBroadcast = errors.New("some error from broadcasting event")
) 

func ProvideUseCaseMocks(t *testing.T) *UseCaseMocks {
	syncOnce.Do(func() {
		repo := repomock.NewMockExtQuerier(t)
		ep := NewMockEventPublisher(t)
		b := NewMockBroadcaster(t)
		ucm = &UseCaseMocks{
			Repository: repo,
			EventPublisher: ep,
			Broadcaster: b,
		}
	})
	return ucm
}