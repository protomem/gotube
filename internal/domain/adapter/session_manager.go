package adapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/infra/flashstore"
)

var _ port.SessionManager = (*SessionManager)(nil)

type SessionManager struct {
	fstore *flashstore.Storage
}

func NewSessionManager(fstore *flashstore.Storage) *SessionManager {
	return &SessionManager{fstore}
}

func (sm *SessionManager) Get(ctx context.Context, token string) (entity.Session, error) {
	var session entity.Session

	if err := sm.fstore.ScanJSON(ctx, sm.buildKey(token), &session); err != nil {
		if errors.Is(err, flashstore.ErrKeyNotFound) {
			return entity.Session{}, entity.ErrSessionNotFound
		}

		return entity.Session{}, err
	}

	return session, nil
}

func (sm *SessionManager) Put(ctx context.Context, session entity.Session) error {
	if err := sm.fstore.PutJSON(ctx, sm.buildKey(session.Token), session, time.Until(session.Expiry)); err != nil {
		return err
	}

	return nil
}

func (sm *SessionManager) Delete(ctx context.Context, token string) error {
	if err := sm.fstore.Delete(ctx, sm.buildKey(token)); err != nil {
		return err
	}

	return nil
}

func (sm *SessionManager) buildKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}
