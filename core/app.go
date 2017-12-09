package core

import (
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type App struct {
	currPos int
	RClient *redis.Client
	chanBytes chan []byte
	chanKeys chan string
	cpMux   sync.RWMutex
	lastReq time.Time
	lrMux   sync.RWMutex
}

func NewApp(client *redis.Client) *App {
	app := &App{
		currPos: 0,
		lastReq: time.Now(),
		RClient: client,
		chanBytes: make(chan []byte, 5000),
		chanKeys: make(chan string, 5000),
	}
	go app.UnmarhalProcess()
	go app.SaveProcess()
	return app
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

func (a *App) IncrCurrPos() {
	defer a.cpMux.Unlock()
	a.cpMux.Lock()
	a.currPos++
}
