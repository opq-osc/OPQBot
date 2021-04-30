package provider

import (
	"container/list"
	"errors"
	"github.com/mcoo/OPQBot/session"
	"sync"
	"time"
)

var pder = &Provider{list: list.New()}

type SessionStore struct {
	qq           int64                       //session id唯一标示
	timeAccessed time.Time                   //最后访问时间
	value        map[interface{}]interface{} //session里面存储的值
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	pder.SessionUpdate(st.qq)
	return nil
}

func (st *SessionStore) Get(key interface{}) (interface{}, error) {
	err := pder.SessionUpdate(st.qq)
	if err != nil {
		return nil, err
	}
	if v, ok := st.value[key]; ok {
		return v, nil
	} else {
		return nil, errors.New("key不存在")
	}
}
func (st *SessionStore) GetString(key interface{}) (string, error) {
	err := pder.SessionUpdate(st.qq)
	if err != nil {
		return "", err
	}
	if v, ok := st.value[key]; ok {
		if v1, ok := v.(string); ok {
			return v1, nil
		} else {
			return "", errors.New("类型转换失败")
		}
	} else {
		return "", errors.New("key不存在")
	}
}
func (st *SessionStore) GetInt(key interface{}) (int, error) {
	err := pder.SessionUpdate(st.qq)
	if err != nil {
		return -1, err
	}
	if v, ok := st.value[key]; ok {
		if v1, ok := v.(int); ok {
			return v1, nil
		} else {
			return -1, errors.New("类型转换失败")
		}
	} else {
		return -1, errors.New("key不存在")
	}
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	pder.SessionUpdate(st.qq)
	return nil
}

func (st *SessionStore) QQ() int64 {
	return st.qq
}

type Provider struct {
	lock     sync.Mutex              //用来锁
	sessions map[int64]*list.Element //用来存储在内存
	list     *list.List              //用来做gc
}

func (pder *Provider) SessionInit(qq int64) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{qq: qq, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[qq] = element
	return newsess, nil
}

func (pder *Provider) SessionRead(qq int64) (session.Session, error) {
	if element, ok := pder.sessions[qq]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := pder.SessionInit(qq)
		return sess, err
	}
}

func (pder *Provider) SessionDestroy(qq int64) error {
	if element, ok := pder.sessions[qq]; ok {
		delete(pder.sessions, qq)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *Provider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).qq)
		} else {
			break
		}
	}
}

func (pder *Provider) SessionUpdate(qq int64) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[qq]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

func init() {
	pder.sessions = make(map[int64]*list.Element, 0)
	session.Register("qq", pder)
}
