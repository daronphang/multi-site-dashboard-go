package usecase

import (
	"context"
	"multi-site-dashboard-go/internal/domain"
	repo "multi-site-dashboard-go/internal/repository"
)

// Each method specifies the topics/queues to publish to.
type StreamPublisher interface {
	PublishToMachineResourceUsage(ctx context.Context, arg domain.CreateMachineResourceUsageParams) error
}

type UseCaseService struct {
	Repository repo.ExtQuerier
	MessagePublisher interface{}
	StreamPublisher StreamPublisher
}

func NewUseCaseService(repo repo.ExtQuerier, sp StreamPublisher) *UseCaseService {
	return &UseCaseService{Repository: repo, StreamPublisher: sp}
}

func (uc *UseCaseService) PublishMsgToWebSocket() {

}