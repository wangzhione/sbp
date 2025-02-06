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
}

func (n *node) size() int {
	if n == nil {
		return 0
	}
	return len(n.persistent)
}

// withNodeFromMaps kvstore 对象转移到 node 上, 转移成功后 node != nil && kvstore 不再可用
func withNodeFromMaps(ctx context.Context, persistent kvstore) context.Context {
	if persistent.size() == 0 {
		return ctx
	}

	nd := new(node)
	// make slices together to reduce malloc cost
	nd.persistent = make([]kv, persistent.size())

	i := 0
	for k, v := range persistent {
		nd.persistent[i].key, nd.persistent[i].val = k, v
		i++
	}

	// kvstore 对象生命周期转给 node, 前者不再可用
	persistent.recycle()

	return withNode(ctx, nd)
}

func (n *node) transferForward() *node {
	return &node{
		persistent: n.persistent,
	}
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

func (n *node) delPersistent(k string) *node {
	if res, ok := remove(n.persistent, k); ok {
		return &node{
			persistent: res,
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
	// need 需要 n.size() != 0
	return context.WithValue(ctx, ctxkey, n)
}

func appendEx(arr []kv, x kv) (res []kv) {
	res = make([]kv, len(arr)+1)
	copy(res, arr)
	res[len(arr)] = x
	return
}
