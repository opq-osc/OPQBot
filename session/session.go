package session

import (
	"fmt"
	"sync"
	"time"
)

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	GetString(key interface{}) string
	GetInt(key interface{}) int
	Delete(key interface{}) error
	QQ() int64
}

type Provider interface {
	SessionInit(qq int64) (Session, error)
	SessionRead(qq int64) (Session, error)
	SessionDestroy(qq int64) error
	SessionGC(maxLifetime int64)
}

var provides = make(map[string]Provider)

func Register(name string, provide Provider) {
	if provide == nil {
		panic("session: Register provide is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provide " + name)
	}
	provides[name] = provide
}

type Manager struct {
	lock        sync.Mutex // protects session
	provider    Provider
	maxlifetime int64
}

func NewManager(provideName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, maxlifetime: maxlifetime}, nil
}

// SessionStart get Session
func (manager *Manager) SessionStart(qq int64) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	session, _ = manager.provider.SessionRead(qq)
	return
}

// SessionDestroy Destroy session
func (manager *Manager) SessionDestroy(qq int64) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionDestroy(qq)
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxlifetime)
	time.AfterFunc(time.Duration(manager.maxlifetime)*time.Second, func() { manager.GC() })
}
