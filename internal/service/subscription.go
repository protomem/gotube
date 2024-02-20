package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ Subscription = (*SubscriptionImpl)(nil)

type SubscriptionDTO struct {
	FromUserNickname string
	ToUserNickname   string
}

type (
	Subscription interface {
		CountSubscribers(ctx context.Context, userNickname string) (int64, error)
		Subscribe(ctx context.Context, dto SubscriptionDTO) error
		Unsubscribe(ctx context.Context, dto SubscriptionDTO) error
	}

	SubscriptionImpl struct {
		repo     repository.Subscription
		userServ User
	}
)

func NewSubscription(repo repository.Subscription, userServ User) *SubscriptionImpl {
	return &SubscriptionImpl{
		repo:     repo,
		userServ: userServ,
	}
}

func (s *SubscriptionImpl) CountSubscribers(ctx context.Context, userNickname string) (int64, error) {
	const op = "service.Subscription.CountSubscribers"

	user, err := s.userServ.GetByNickname(ctx, userNickname)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	count, err := s.repo.CountByToUserID(ctx, user.ID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (s *SubscriptionImpl) Subscribe(ctx context.Context, dto SubscriptionDTO) error {
	const op = "service.Subscription.Subscribe"

	fromUser, toUser, err := s.getUsers(ctx, dto.FromUserNickname, dto.ToUserNickname)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := s.repo.Create(ctx, repository.CreateSubscriptionDTO{
		FromUserID: fromUser.ID,
		ToUserID:   toUser.ID,
	}); err != nil {
		if errors.Is(err, model.ErrSubscriptionExists) {
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SubscriptionImpl) Unsubscribe(ctx context.Context, dto SubscriptionDTO) error {
	const op = "service.Subscription.Unsubscribe"

	fromUser, toUser, err := s.getUsers(ctx, dto.FromUserNickname, dto.ToUserNickname)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	sub, err := s.repo.GetByFromUserIDAndToUserID(ctx, fromUser.ID, toUser.ID)
	if err != nil {
		if errors.Is(err, model.ErrSubscriptionNotFound) {
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.repo.Delete(ctx, sub.ID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SubscriptionImpl) getUsers(ctx context.Context, fromUserNickname, toUserNickname string) (model.User, model.User, error) {
	fromUser, err := s.userServ.GetByNickname(ctx, fromUserNickname)
	if err != nil {
		return model.User{}, model.User{}, fmt.Errorf("get toUser: %w", err)
	}

	toUser, err := s.userServ.GetByNickname(ctx, toUserNickname)
	if err != nil {
		return model.User{}, model.User{}, fmt.Errorf("get fromUser: %w", err)
	}

	return fromUser, toUser, nil
}
