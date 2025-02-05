package metainfo

import (
	"context"
)

type kv struct {
	key string
	val string
}

type node struct {
	persistent []kv
	transient  []kv
	stale      []kv
}

func (n *node) size() int {
	if n == nil {
		return 0
	}
	return len(n.persistent) + len(n.transient) + len(n.stale)
}

// withNodeFromMaps kvstore 对象转移到 node 上, 转移成功后 node != nil && kvstore 不再可用
func withNodeFromMaps(ctx context.Context, persistent, transient, stale kvstore) context.Context {
	// need 需要 ps + ts + sz > 0
	ps, ts, sz := persistent.size(), transient.size(), stale.size()
	if ps+ts+sz == 0 {
		return ctx
	}

	nd := new(node)
	// make slices together to reduce malloc cost
	kvs := make([]kv, ps+ts+sz)
	if ps > 0 {
		nd.persistent = kvs[:ps]
	}
	if ts > 0 {
		nd.transient = kvs[ps : ps+ts]
	}
	if ts > 0 {
		nd.stale = kvs[ps+ts:]
	}

	i := 0
	for k, v := range persistent {
		nd.persistent[i].key, nd.persistent[i].val = k, v
		i++
	}
	i = 0
	for k, v := range transient {
		nd.transient[i].key, nd.transient[i].val = k, v
		i++
	}
	i = 0
	for k, v := range stale {
		nd.stale[i].key, nd.stale[i].val = k, v
		i++
	}

	// kvstore 对象生命周期转给 node, 前者不再可用
	persistent.recycle()
	transient.recycle()
	stale.recycle()

	return withNode(ctx, nd)
}

func (n *node) transferForward() *node {
	return &node{
		persistent: n.persistent,
		stale:      n.transient,
	}
}

func (n *node) addTransient(k, v string) *node {
	if res, ok := remove(n.stale, k); ok {
		return &node{
			persistent: n.persistent,
			transient: appendEx(n.transient, kv{
				key: k,
				val: v,
			}),
			stale: res,
		}
	}

	if idx, ok := search(n.transient, k); ok {
		if n.transient[idx].val == v {
			return n
		}
		r := *n
		r.transient = make([]kv, len(n.transient))
		copy(r.transient, n.transient)
		r.transient[idx].val = v
		return &r
	}

	r := *n
	r.transient = appendEx(r.transient, kv{
		key: k,
		val: v,
	})
	return &r
}

func (n *node) addPersistent(k, v string) *node {
	if idx, ok := search(n.persistent, k); ok {
		if n.persistent[idx].val == v {
			return n
		}
		r := *n
		r.persistent = make([]kv, len(n.persistent))
		copy(r.persistent, n.persistent)
		r.persistent[idx].val = v
		return &r
	}
	r := *n
	r.persistent = appendEx(r.persistent, kv{
		key: k,
		val: v,
	})
	return &r
}

func (n *node) delTransient(k string) (r *node) {
	if res, ok := remove(n.stale, k); ok {
		return &node{
			persistent: n.persistent,
			transient:  n.transient,
			stale:      res,
		}
	}
	if res, ok := remove(n.transient, k); ok {
		return &node{
			persistent: n.persistent,
			transient:  res,
			stale:      n.stale,
		}
	}
	return n
}

func (n *node) delPersistent(k string) (r *node) {
	if res, ok := remove(n.persistent, k); ok {
		return &node{
			persistent: res,
			transient:  n.transient,
			stale:      n.stale,
		}
	}
	return n
}

func search(kvs []kv, key string) (idx int, ok bool) {
	for i := range kvs {
		if kvs[i].key == key {
			return i, true
		}
	}
	return
}

func remove(kvs []kv, key string) (res []kv, removed bool) {
	if idx, ok := search(kvs, key); ok {
		if cnt := len(kvs); cnt == 1 {
			removed = true
			return
		}
		res = make([]kv, len(kvs)-1)
		copy(res, kvs[:idx])
		copy(res[idx:], kvs[idx+1:])
		return res, true
	}
	return kvs, false
}

type ctxkeytype struct{}

var ctxkey ctxkeytype

func getNode(ctx context.Context) *node {
	val, _ := ctx.Value(ctxkey).(*node)
	return val
}

func withNode(ctx context.Context, n *node) context.Context {
	// return original ctx if no invalid key in map
	if n.size() == 0 {
		return ctx
	}
	return context.WithValue(ctx, ctxkey, n)
}

func appendEx(arr []kv, x kv) (res []kv) {
	res = make([]kv, len(arr)+1)
	copy(res, arr)
	res[len(arr)] = x
	return
}
