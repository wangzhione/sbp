package metainfo

import (
	"context"
)

// The prefix listed below may be used to tag the types of values when there is no context to carry them.
// HTTP header prefixes.
const (
	PrefixPersistent = "Rpc-Persist-" // to textproto.CanonicalMIMEHeaderKey
	lenPP            = len(PrefixPersistent)
)

// **Using empty string as key or value is not support.**

// GetPersistentValue retrieves the persistent value set into the context by the given key.
func GetPersistentValue(ctx context.Context, k string) (v string, ok bool) {
	if n := getNode(ctx); n.size() > 0 {
		if i := search(n.persistent, k); i >= 0 {
			return n.persistent[i].val, true
		}
	}
	return
}

// GetAllPersistentValues retrieves all persistent values.
func GetAllPersistentValues(ctx context.Context) (m map[string]string) {
	if n := getNode(ctx); n.size() > 0 {
		m = make(map[string]string, len(n.persistent))
		for _, kv := range n.persistent {
			m[kv.key] = kv.val
		}
	}
	return
}

// RangePersistentValues calls fn sequentially for each persistent kv.
// If fn returns false, range stops the iteration.
func RangePersistentValues(ctx context.Context, fn func(k, v string) bool) {
	if n := getNode(ctx); n.size() > 0 {
		for _, kv := range n.persistent {
			if !fn(kv.key, kv.val) {
				break
			}
		}
	}
}

// WithPersistentValue sets the value into the context by the given key.
// This value will be propagated to the services along the RPC call chain.
func WithPersistentValue(ctx context.Context, k, v string) context.Context {
	if len(k) == 0 || len(v) == 0 {
		return ctx
	}

	if n := getNode(ctx); n != nil {
		if m := n.addPersistent(k, v); m != n {
			return withNode(ctx, m)
		}
		return ctx
	}
	return withNode(ctx, &node{
		persistent: []kv{{key: k, val: v}},
	})
}

// DelPersistentValue deletes a persistent key/value from the current context.
// Since empty string value is not valid, we could just set the value to be empty.
func DelPersistentValue(ctx context.Context, k string) context.Context {
	if len(k) == 0 {
		return ctx
	}
	if n := getNode(ctx); n.size() > 0 {
		if m := n.delPersistent(k); m != n {
			return withNode(ctx, m)
		}
	}
	return ctx
}

// CountPersistentValues counts the length of persisten KV pairs
func CountPersistentValues(ctx context.Context) int {
	return getNode(ctx).size()
}

// WithPersistentValues sets the values into the context by the given keys.
// This value will be propagated to the services along the RPC call chain.
func WithPersistentValues(ctx context.Context, kvs ...string) context.Context {
	if len(kvs) == 0 || len(kvs)%2 != 0 {
		return ctx
	}

	kvLen := len(kvs) / 2

	n := &node{}
	if m := getNode(ctx); m.size() > 0 {
		n.persistent = make([]kv, len(m.persistent), len(m.persistent)+kvLen)
		copy(n.persistent, m.persistent)
	} else {
		n.persistent = make([]kv, 0, kvLen)
	}

	for i := 0; i < kvLen; i++ {
		key, val := kvs[i*2], kvs[i*2+1]
		if len(key) == 0 || len(val) == 0 {
			continue
		}

		if idx := search(n.persistent, key); idx >= 0 {
			if n.persistent[idx].val != val {
				n.persistent[idx].val = val
			}
			continue
		}

		n.persistent = append(n.persistent, kv{key: key, val: val})
	}

	return withNode(ctx, n)
}
