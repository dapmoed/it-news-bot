package sessions

import (
	"errors"
	"it-news-bot/internal/chains"
	"time"
)

const (
	defaultExpiredDuration = 1 * time.Minute
)

var (
	ErrNotFound = errors.New("not found")
)

type StorageSession struct {
	sessions map[int64]Session
}

type Session struct {
	userId  int64
	chain   *chains.Chain
	expired time.Time
}

func (s *Session) GetChain() *chains.Chain {
	return s.chain
}

func New() *StorageSession {
	return &StorageSession{
		sessions: make(map[int64]Session),
	}
}

func (s *StorageSession) Get(userId int64) (Session, error) {
	s.Clear()
	if val, ok := s.sessions[userId]; ok {
		if time.Now().Before(s.sessions[userId].expired) {
			return val, nil
		}
		return Session{}, ErrNotFound
	}
	return Session{}, ErrNotFound
}

func (s *StorageSession) Add(userId int64, chain *chains.Chain) {
	s.sessions[userId] = Session{
		userId:  userId,
		chain:   chain,
		expired: time.Now().Add(defaultExpiredDuration),
	}
}

func (s *StorageSession) Clear() {
	for i, v := range s.sessions {
		if v.expired.Before(time.Now()) {
			delete(s.sessions, i)
		}
	}
}
