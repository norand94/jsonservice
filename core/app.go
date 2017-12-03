package core

import (
	"time"
	"sync"
	"github.com/go-redis/redis"
)

type App struct {
	currPos int
	RClient *redis.Client
	cpMux   sync.RWMutex
	lastReq time.Time
	lrMux   sync.RWMutex
}

func NewApp(client *redis.Client) *App  {
	return &App{
		currPos: 0,
		lastReq: time.Now(),
		RClient: client,
	}
}

func (a *App) LastReq() time.Time {
	defer a.lrMux.RUnlock()
	a.lrMux.RLock()
	return a.lastReq
}

func (a *App) SetLastReq(t time.Time) {
	defer a.lrMux.Unlock()
	a.lrMux.Lock()
	a.lastReq = t
}

func (a *App) CurrPos() int {
	defer a.cpMux.RUnlock()
	a.cpMux.RLock()
	return a.currPos
}

func (a *App)IncrCurrPos () {
	defer a.cpMux.Unlock()
	a.cpMux.Lock()
	a.currPos++
}