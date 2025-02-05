package metainfo

import (
	"context"
	"sync"
)

type bwCtxKeyType int

const (
	bwCtxKeySend bwCtxKeyType = iota
	bwCtxKeyRecv
)

type bwCtxValue struct {
	sync.RWMutex
	kvs map[string]string
}

func newBackwardCtxValues() *bwCtxValue {
	return &bwCtxValue{
		kvs: make(map[string]string),
	}
}

func (p *bwCtxValue) get(k string) (v string, ok bool) {
	p.RLock()
	defer p.RUnlock()

	v, ok = p.kvs[k]
	return
}

func (p *bwCtxValue) getAll() (m map[string]string) {
	p.RLock()
	defer p.RUnlock()

	if len(p.kvs) > 0 {
		m = make(map[string]string, len(p.kvs))
		for k, v := range p.kvs {
			m[k] = v
		}
	}
	return
}

func (p *bwCtxValue) set(k, v string) {
	p.Lock()
	defer p.Unlock()

	p.kvs[k] = v
}

// setMany len(kvs) 必须是偶数
func (p *bwCtxValue) setMany(kvs []string) bool {
	if len(kvs) == 0 {
		return true
	}
	if len(kvs)%2 != 0 {
		return false
	}

	p.Lock()
	for i := 0; i < len(kvs); i += 2 {
		p.kvs[kvs[i]] = kvs[i+1]
	}
	p.Unlock()
	return true
}

func (p *bwCtxValue) setMap(kvs map[string]string) {
	if len(kvs) == 0 {
		return
	}

	p.Lock()
	for k, v := range kvs {
		p.kvs[k] = v
	}
	p.Unlock()
}

// WithBackwardValues returns a new context that allows passing key-value pairs
// backward with `SetBackwardValue` from any derived context.
func WithBackwardValues(ctx context.Context) context.Context {
	if _, ok := ctx.Value(bwCtxKeyRecv).(*bwCtxValue); ok {
		return ctx
	}
	return context.WithValue(ctx, bwCtxKeyRecv, newBackwardCtxValues())
}

// RecvBackwardValue gets a value associated with the given key that is set by
// `SetBackwardValue` or `SetBackwardValues`.
func RecvBackwardValue(ctx context.Context, key string) (val string, ok bool) {
	if p, exist := ctx.Value(bwCtxKeyRecv).(*bwCtxValue); exist {
		val, ok = p.get(key)
	}
	return
}

// RecvAllBackwardValues is the batched version of RecvBackwardValue.
func RecvAllBackwardValues(ctx context.Context) (m map[string]string) {
	if p, ok := ctx.Value(bwCtxKeyRecv).(*bwCtxValue); ok {
		return p.getAll()
	}
	return
}

// SetBackwardValue sets a key value pair into the context.
func SetBackwardValue(ctx context.Context, key, val string) bool {
	if p, ok := ctx.Value(bwCtxKeyRecv).(*bwCtxValue); ok {
		p.set(key, val)
		return true
	}
	return false
}

// SetBackwardValues is the batched version of `SetBackwardValue`.
func SetBackwardValues(ctx context.Context, kvs ...string) bool {
	if p, ok := ctx.Value(bwCtxKeyRecv).(*bwCtxValue); ok {
		return p.setMany(kvs)
	}
	return false
}

// SetBackwardValuesFromMap is the batched version of `SetBackwardValue`.
func SetBackwardValuesFromMap(ctx context.Context, kvs map[string]string) bool {
	if p, ok := ctx.Value(bwCtxKeyRecv).(*bwCtxValue); ok {
		p.setMap(kvs)
		return true
	}
	return false
}

// WithBackwardValuesToSend returns a new context that collects key-value
// pairs set with `SendBackwardValue` or `SendBackwardValues` into any
// derived context.
func WithBackwardValuesToSend(ctx context.Context) context.Context {
	if _, ok := ctx.Value(bwCtxKeySend).(*bwCtxValue); ok {
		return ctx
	}
	return context.WithValue(ctx, bwCtxKeySend, newBackwardCtxValues())
}

// SendBackwardValue sets a key-value pair into the context for sending to
// a remote endpoint.
// Note that the values can not be retrieved with `RecvBackwardValue` from
// the same context.
func SendBackwardValue(ctx context.Context, key, val string) bool {
	if p, ok := ctx.Value(bwCtxKeySend).(*bwCtxValue); ok {
		p.set(key, val)
		return true
	}
	return false
}

// SendBackwardValues is the batched version of `SendBackwardValue`.
func SendBackwardValues(ctx context.Context, kvs ...string) bool {
	if p, ok := ctx.Value(bwCtxKeySend).(*bwCtxValue); ok {
		return p.setMany(kvs)
	}
	return false
}

// SendBackwardValuesFromMap is the batched version of `SendBackwardValue`.
func SendBackwardValuesFromMap(ctx context.Context, kvs map[string]string) bool {
	if p, ok := ctx.Value(bwCtxKeySend).(*bwCtxValue); ok {
		p.setMap(kvs)
		return true
	}
	return false
}

// GetBackwardValueToSend gets a value associated with the given key that is set by
// `SendBackwardValue`, `SendBackwardValues` or `SendBackwardValuesFromMap`.
func GetBackwardValueToSend(ctx context.Context, key string) (val string, ok bool) {
	if p, exist := ctx.Value(bwCtxKeySend).(*bwCtxValue); exist {
		return p.get(key)
	}
	return
}

// AllBackwardValuesToSend retrieves all key-values pairs set by `SendBackwardValue`
// or `SendBackwardValues` from the given context.
// This function is designed for frameworks, common developers should not use it.
func AllBackwardValuesToSend(ctx context.Context) (m map[string]string) {
	if p, ok := ctx.Value(bwCtxKeySend).(*bwCtxValue); ok {
		return p.getAll()
	}
	return
}
