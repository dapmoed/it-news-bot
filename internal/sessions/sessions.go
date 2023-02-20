package sessions

import (
	"errors"
	"it-news-bot/internal/chains"
	"time"
)

const (
	defaultExpiredDuration = 10 * time.Second
)

var (
	ErrNotFound = errors.New("not found")
)

type StorageSession struct {
	sessions map[int64]*Session
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
	strgSess := &StorageSession{
		sessions: make(map[int64]*Session),
	}
	go strgSess.RunCollector()
	return strgSess
}

func (s *StorageSession) Get(userId int64) (*Session, error) {
	s.Clear()
	if val, ok := s.sessions[userId]; ok {
		if time.Now().Before(s.sessions[userId].expired) {
			return val, nil
		}
		return nil, ErrNotFound
	}
	return nil, ErrNotFound
}

func (s *Session) Extend() {
	s.expired = time.Now().Add(defaultExpiredDuration)
}

func (s *StorageSession) Add(userId int64, chain *chains.Chain) {
	s.sessions[userId] = &Session{
		userId:  userId,
		chain:   chain,
		expired: time.Now().Add(defaultExpiredDuration),
	}
}

func (s *StorageSession) Clear() {
	for i, v := range s.sessions {
		if v.expired.Before(time.Now()) {
			delete(s.sessions, i)
			continue
		}
		if v.chain.IsEnded() {
			delete(s.sessions, i)
			continue
		}
	}
}

func (s *StorageSession) RunCollector() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			s.Clear()
		}
	}
}
