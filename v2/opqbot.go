package OPQBot

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/opq-osc/OPQBot/v2/errors"
	"github.com/opq-osc/OPQBot/v2/events"
	"github.com/rotisserie/eris"
	"net/url"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"time"
)

type Core struct {
	ApiUrl      *url.URL
	events      map[events.EventName][]events.EventCallbackFunc
	lock        sync.RWMutex
	err         error
	client      *websocket.Conn
	handlePanic func(any)

	retryCount, MaxRetryCount int

	done chan struct{}
}

func (c *Core) HandlePanic(h func(any)) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.handlePanic = h
}

func (c *Core) On(event events.EventName, callback events.EventCallbackFunc) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.events[event] = append(c.events[event], callback)
}

func (c *Core) closeEvent() {
	log.Info("即将关闭")
}

func (c *Core) ListenAndWait(ctx context.Context) (e error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	go func() {
		select {
		case <-interrupt:
			log.Info("用户关闭程序")
			cancel()
			if c.client != nil {
				c.client.Close()
			}

		case <-ctx.Done():
		}
	}()

	c.done = make(chan struct{}, 1)
	defer func() {
		log.Debug(e)
		if e != errors.ErrorContextCanceled {
			c.retryCount++
			if c.retryCount > c.MaxRetryCount {
				log.Info("超出最大重连次数")
				return
			}
			log.Warnf("连接出错，将进行第%d重连操作,按Ctrl+C取消重试", c.retryCount)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(c.retryCount) * time.Second):
			}
			c.err = nil
			e = c.ListenAndWait(ctx)
			return
		}
		c.closeEvent()
	}()
	var err error
	c.client, _, err = websocket.DefaultDialer.DialContext(ctx, "ws://"+c.ApiUrl.Host+"/ws", nil)
	if err != nil {
		return err
	}
	defer func() {
		if c.client != nil {
			c.client.Close()
		}
	}()
	c.retryCount = 0
	log.Info("连接成功")
	go func() {
		defer close(c.done)
		for {
			_, message, err := c.client.ReadMessage()
			select {
			case <-ctx.Done():
				c.err = errors.ErrorContextCanceled
				return
			default:
			}
			if err != nil {
				c.err = err
				return
			}
			event, err := events.New(c.ApiUrl.Scheme+"://"+c.ApiUrl.Host, message)
			if err != nil {
				log.Error("error:", eris.ToString(err, true))
				continue
			}
			log.Debug(string(message))
			var callbacks []events.EventCallbackFunc
			c.lock.RLock()
			callbacks = c.events[event.GetEventName()]
			c.lock.RUnlock()
			go func() {
				defer func() {
					if err := recover(); err != nil {
						if c.handlePanic != nil {
							c.handlePanic(err)
						} else {
							log.Infof("event handle function panic: %s \n%s", err, string(debug.Stack()))
						}
					}
				}()
				for _, v := range callbacks {
					v(ctx, event)
				}
			}()

		}

	}()

	<-c.done

	return c.err
}

func NewCore(api string, maxRetryCount int) (*Core, error) {
	u, err := url.Parse(api)
	if err != nil {
		return nil, err
	}
	return &Core{
		ApiUrl:        u,
		events:        make(map[events.EventName][]events.EventCallbackFunc),
		lock:          sync.RWMutex{},
		done:          nil,
		MaxRetryCount: maxRetryCount,
	}, nil
}
