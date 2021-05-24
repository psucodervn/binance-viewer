package runner

import (
	"sync"

	"github.com/rs/zerolog/log"

	"copytrader/internal/model"
)

type Manager struct {
	db      *model.Database
	runners map[string]*AccountRunner

	mu sync.RWMutex
}

func NewManager(db *model.Database) *Manager {
	return &Manager{db: db, runners: map[string]*AccountRunner{}}
}

func (m *Manager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	users := m.db.Users()
	for _, u := range users {
		for _, acc := range u.Accounts {
			m.runners[acc.ApiKey] = NewAccountRunner(acc.ApiKey, acc.SecretKey)
		}
	}

	for _, r := range m.runners {
		if err := r.OnUpdate(); err != nil {
			log.Err(err).Send()
		}
	}
}
