package usecase

import (
	"context"
	"multi-site-dashboard-go/internal/domain"
	repo "multi-site-dashboard-go/internal/repository"
)

// Each method specifies the topics/queues to publish to.
type EventPublisher interface {
	PublishDataToMachineResourceUsage(ctx context.Context, arg domain.CreateMachineResourceUsageParams) error
}

type Broadcaster interface {
	Broadcast(ctx context.Context, data []byte) error
}

type UseCaseService struct {
	Repository repo.ExtQuerier
	EventPublisher EventPublisher
	Broadcaster Broadcaster
}

func NewUseCaseService(repo repo.ExtQuerier, ep EventPublisher, b Broadcaster) *UseCaseService {
	return &UseCaseService{Repository: repo, EventPublisher: ep, Broadcaster: b}
}

func (uc *UseCaseService) PublishMsgToWebSocket() {

}