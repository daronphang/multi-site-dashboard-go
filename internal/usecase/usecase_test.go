package usecase

import (
	"errors"
	"multi-site-dashboard-go/internal/repository"
	"sync"
	"testing"
)

type useCaseMocks struct {
	repository *repository.MockExtQuerier
	eventPublisher *MockEventPublisher
	broadcaster *MockBroadcaster
	usecase *UseCaseService
}

var (
	ucm *useCaseMocks
	syncOnce sync.Once
	errDB = errors.New("some error from db")
	errEventPublisher = errors.New("some error from publishing event")
	errBroadcast = errors.New("some error from broadcasting event")
) 

func provideUseCaseMocks(t *testing.T) *useCaseMocks {
	syncOnce.Do(func() {
		repo := repository.NewMockExtQuerier(t)
		ep := NewMockEventPublisher(t)
		b := NewMockBroadcaster(t)
		uc := NewUseCaseService(repo, ep, b)
		ucm = &useCaseMocks{
			repository: repo,
			eventPublisher: ep,
			broadcaster: b,
			usecase: uc,
		}
	})
	return ucm
}