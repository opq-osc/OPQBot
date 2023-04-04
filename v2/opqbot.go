package OPQBot

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/opq-osc/OPQBot/v2/errors"
	"github.com/opq-osc/OPQBot/v2/events"
	"log"
	"net/url"
	"sync"
)

type Core struct {
	ApiUrl *url.URL
	events map[events.EventName][]events.EventCallbackFunc
	lock   sync.RWMutex

	done chan struct{}
}

func (c *Core) On(event events.EventName, callback events.EventCallbackFunc) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.events[event] = append(c.events[event], callback)
}

func (c *Core) ListenAndWait(ctx context.Context) error {
	c.done = make(chan struct{}, 1)
	client, _, err := websocket.DefaultDialer.Dial("ws://"+c.ApiUrl.Host+"/ws", nil)
	if err != nil {
		return err
	}
	defer client.Close()
	go func() {
		defer close(c.done)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			_, message, err := client.ReadMessage()
			if err != nil {
				log.Println("error:", err)
				return
			}
			event, err := events.New(c.ApiUrl.Scheme+"://"+c.ApiUrl.Host+"/v1/LuaApiCaller", message)
			if err != nil {
				log.Println("error:", err)
				continue
			}
			var callbacks []events.EventCallbackFunc
			c.lock.RLock()
			callbacks = c.events[event.GetEventName()]
			c.lock.RUnlock()
			for _, v := range callbacks {
				v(ctx, event)
			}
		}
	}()

	select {
	case <-c.done:
		// 正常退出
		return nil

	case <-ctx.Done():
		return errors.ErrorContextCanceled
	}
}

func NewCore(api string) (*Core, error) {
	u, err := url.Parse(api)
	if err != nil {
		return nil, err
	}
	return &Core{
		ApiUrl: u,
		events: make(map[events.EventName][]events.EventCallbackFunc),
		lock:   sync.RWMutex{},
		done:   nil,
	}, nil
}
