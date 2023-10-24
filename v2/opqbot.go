package OPQBot

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/jasonlvhit/gocron"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/errors"
	"github.com/opq-osc/OPQBot/v2/events"
	"github.com/rotisserie/eris"
)

type Core struct {
	ApiUrl                    *url.URL
	events                    map[events.EventName][]events.EventCallbackFunc
	lock                      sync.RWMutex
	err                       error
	client                    *websocket.Conn
	handlePanic               func(any)
	retryCount, MaxRetryCount int
	autoSignToken             bool
	apibase                   string
	botQQ                     *int64
	groupQQ                   *int64

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

func (c *Core) ListenAndWait(ctx context.Context) (e error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	go c.MakeAutoSign()
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
			log.Debug(string(message))
			event, err := events.New(message)
			if err != nil {
				log.Error("error:", eris.ToString(err, true))
				continue
			}

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

// MakeAutoSign 定时任务 自动签到
func (c *Core) MakeAutoSign() {
	if c.autoSignToken && c.botQQ != nil && c.groupQQ != nil {
		log.Info("已经开启自动签到，每天1点在官方群发送签到信息")
		s := gocron.NewScheduler()
		err := s.Every(1).Days().At("01:00").Do(func() {
			log.Info("执行自动签到任务")
			apiBuilder.New(c.apibase, *c.botQQ).SendMsg().GroupMsg().ToUin(*c.groupQQ).TextMsg("签到").Do(context.Background())
		})
		if err != nil {
			return
		}
		<-s.Start()
	}
}

type CoreOpt func(*Core)

func WithMaxRetryCount(maxCount int) CoreOpt {
	return func(c *Core) {
		c.MaxRetryCount = maxCount
	}
}

func WithAutoSignToken(botQQ int64, groupQQ int64) CoreOpt {
	return func(c *Core) {
		c.autoSignToken = true
		c.botQQ = &botQQ
		c.groupQQ = &groupQQ
	}
}

func NewCore(api string, opt ...CoreOpt) (*Core, error) {
	u, err := url.Parse(api)
	if err != nil {
		return nil, err
	}
	c := &Core{
		ApiUrl:        u,
		apibase:       api,
		events:        make(map[events.EventName][]events.EventCallbackFunc),
		lock:          sync.RWMutex{},
		done:          nil,
		MaxRetryCount: 10,
		autoSignToken: false,
	}
	for _, v := range opt {
		v(c)
	}
	return c, nil
}
