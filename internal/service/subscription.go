package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ Subscription = (*SubscriptionImpl)(nil)

type (
	SubscribeDTO struct {
		FromUserID uuid.UUID
		ToUserID   uuid.UUID
	}

	UnsubscribeDTO struct {
		FromUserID uuid.UUID
		ToUserID   uuid.UUID
	}

	Subscription interface {
		FindByFromUserNickname(ctx context.Context, fromUserNickname string) ([]model.Subscription, error)

		Subscribe(ctx context.Context, dto SubscribeDTO) error
		Unsubscribe(ctx context.Context, dto UnsubscribeDTO) error
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

func (serv *SubscriptionImpl) FindByFromUserNickname(
	ctx context.Context,
	fromUserNickname string,
) ([]model.Subscription, error) {
	const op = "service.Subscription.FindByFromUserNickname"
	var err error

	fromUser, err := serv.userServ.GetByNickname(ctx, fromUserNickname)
	if err != nil {
		return []model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	subs, err := serv.repo.FindByFromUserID(ctx, fromUser.ID)
	if err != nil {
		return []model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	return subs, nil
}

func (serv *SubscriptionImpl) Subscribe(ctx context.Context, dto SubscribeDTO) error {
	const op = "service.Subscription.Subscribe"
	var err error

	_, err = serv.repo.GetByFromUserAndToUser(ctx, dto.FromUserID, dto.ToUserID)
	if err != nil && !errors.Is(err, model.ErrSubscriptionNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = serv.repo.Create(ctx, repository.CreateSubscriptionDTO(dto))
	if err != nil {
		if errors.Is(err, model.ErrSubscriptionExists) {
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (serv *SubscriptionImpl) Unsubscribe(ctx context.Context, dto UnsubscribeDTO) error {
	const op = "service.Subscription.Unsubscribe"
	var err error

	sub, err := serv.repo.GetByFromUserAndToUser(ctx, dto.FromUserID, dto.ToUserID)
	if err != nil {
		if errors.Is(err, model.ErrSubscriptionNotFound) {
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	err = serv.repo.Delete(ctx, sub.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}