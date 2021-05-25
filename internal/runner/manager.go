package runner

import (
	"github.com/adshao/go-binance/v2/futures"
	"sync"

	"github.com/rs/zerolog/log"

	"copytrader/internal/model"
)

type Manager struct {
	db      *model.Database
	runners map[string]*AccountRunner
	onEvent OnEvent

	mu sync.RWMutex
}

type OnEvent func(u model.User, acc model.Account, ev *futures.WsUserDataEvent)

func NewManager(db *model.Database) *Manager {
	return &Manager{db: db, runners: map[string]*AccountRunner{}}
}

func (m *Manager) Subscribe(fn OnEvent) {
	m.mu.Lock()
	m.onEvent = fn
	m.mu.Unlock()
}

func (m *Manager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	users := m.db.Users()
	onEvent := m.onEvent
	for i := range users {
		u := users[i]
		for j := range u.Accounts {
			acc := u.Accounts[j]
			fn := func(ev *futures.WsUserDataEvent) {
				if onEvent != nil {
					onEvent(u, acc, ev)
				}
			}
			r := NewAccountRunner(acc.ApiKey, acc.SecretKey)
			if err := r.OnUpdate(fn); err != nil {
				log.Err(err).Send()
				continue
			}
			m.runners[u.ID+":"+acc.ApiKey] = r
		}
	}
}
