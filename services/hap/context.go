package hap

import (
	"net"
	"net/http"
	"sync"
)

type Context interface {
	GetSessionForRequest(r *http.Request) Session
	SetSessionForConnection(s Session, conn net.Conn)
}

type context struct {
	storage map[interface{}]interface{}
	mutex   *sync.Mutex
}

//
func (ctx *context) GetKey(c net.Conn) interface{} {
	return c.RemoteAddr().String()
}

func (ctx *context) GetConnectionKey(r *http.Request) interface{} {
	return r.RemoteAddr
}

func (ctx *context) Get(key interface{}) interface{} {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.storage[key]
}

func (ctx *context) Set(key interface{}, val interface{}) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.storage[key] = val

}

func (ctx *context) Delete(key interface{}) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	delete(ctx.storage, key)
}

//Hap context
func (ctx *context) GetSessionForRequest(r *http.Request) Session {
	key := ctx.GetConnectionKey(r)
	return ctx.Get(key).(Session)
}

func (ctx *context) SetSessionForConnection(s Session, conn net.Conn) {
	key := ctx.GetKey(conn)
	ctx.Set(key, s)
}
