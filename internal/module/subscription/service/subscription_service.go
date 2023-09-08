package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/module/subscription/model"
	"github.com/protomem/gotube/internal/module/subscription/repository"
	userserv "github.com/protomem/gotube/internal/module/user/service"
)

var _ SubscriptionService = (*SubscriptionServiceImpl)(nil)

type SubscribeDTO struct {
	FromUserID     uuid.UUID
	ToUserNickname string
}

type (
	SubscriptionService interface {
		Subscribe(ctx context.Context, dto SubscribeDTO) error
	}

	SubscriptionServiceImpl struct {
		userServ userserv.UserService

		subscriptionRepo repository.SubscriptionRepository
	}
)

func NewSubscriptionService(
	userServ userserv.UserService,
	subscriptionRepo repository.SubscriptionRepository,
) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{
		userServ:         userServ,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *SubscriptionServiceImpl) Subscribe(ctx context.Context, dto SubscribeDTO) error {
	const op = "SubscriptionService.CreateSubscription"
	var err error

	toUser, err := s.userServ.FindOneUserByNickname(ctx, dto.ToUserNickname)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.subscriptionRepo.CreateSubscription(ctx, repository.CreateSubscriptionDTO{
		FromUserID: dto.FromUserID,
		ToUserID:   toUser.ID,
	})
	if err != nil && !errors.Is(err, model.ErrSubscriptionAlreadyExists) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
