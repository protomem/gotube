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

type UnsubscribeDTO struct {
	FromUserID     uuid.UUID
	ToUserNickname string
}

type StatisticsDTO struct {
	CountSubscriptions uint64
	CountSubscribers   uint64
}

type (
	SubscriptionService interface {
		FindAllSubscriptionsByFromUserNickname(ctx context.Context, fromUserNickname string) ([]model.Subscription, error)

		Subscribe(ctx context.Context, dto SubscribeDTO) error
		Unsubscribe(ctx context.Context, dto UnsubscribeDTO) error

		Statistics(ctx context.Context, userNickname string) (StatisticsDTO, error)
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

func (s *SubscriptionServiceImpl) FindAllSubscriptionsByFromUserNickname(
	ctx context.Context,
	fromUserNickname string,
) ([]model.Subscription, error) {
	const op = "SubscriptionService.FindAllSubscriptionsByFromUserNickname"
	var err error

	fromUser, err := s.userServ.FindOneUserByNickname(ctx, fromUserNickname)
	if err != nil {
		return []model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	subscriptions, err := s.subscriptionRepo.FindAllSubscriptionsByFromUserID(ctx, fromUser.ID)
	if err != nil {
		return []model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	return subscriptions, nil
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

func (s *SubscriptionServiceImpl) Unsubscribe(ctx context.Context, dto UnsubscribeDTO) error {
	const op = "SubscriptionService.DeleteSubscription"
	var err error

	toUser, err := s.userServ.FindOneUserByNickname(ctx, dto.ToUserNickname)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.subscriptionRepo.DeleteSubscription(ctx, repository.DeleteSubscriptionDTO{
		FromUserID: dto.FromUserID,
		ToUserID:   toUser.ID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SubscriptionServiceImpl) Statistics(ctx context.Context, userNickname string) (StatisticsDTO, error) {
	const op = "SubscriptionService.Statistics"
	var err error

	user, err := s.userServ.FindOneUserByNickname(ctx, userNickname)
	if err != nil {
		return StatisticsDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	countSubscriptions, err := s.subscriptionRepo.CountSubscriptionsByFromUserID(ctx, user.ID)
	if err != nil {
		return StatisticsDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	countSubscribers, err := s.subscriptionRepo.CountSubscriptionsByToUserID(ctx, user.ID)
	if err != nil {
		return StatisticsDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	return StatisticsDTO{
		CountSubscriptions: countSubscriptions,
		CountSubscribers:   countSubscribers,
	}, nil
}
