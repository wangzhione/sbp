package metainfo

import (
	"context"
	"strings"
)

// SetMetaInfoFromMap retrieves metainfo key-value pairs from the given map and sets then into the context.
// Only those keys with prefixes defined in this module would be used.
// If the context has been carrying metanifo pairs, they will be merged as a basis.
func SetMetaInfoFromMap(ctx context.Context, m map[string]string) context.Context {
	// need ctx != nil, 永远别用 nil context 去玩
	if len(m) == 0 {
		return ctx
	}

	nd := getNode(ctx)
	if nd.size() == 0 {
		// fast path
		return newCtxFromMap(ctx, m)
	}

	// inherit from node
	persistent := newkvtostore(nd.persistent)

	// insert new kvs from m to node
	for k, v := range m {
		if len(k) == 0 || len(v) == 0 {
			continue
		}
		switch {
		case strings.HasPrefix(k, PrefixPersistent):
			if len(k) > lenPP {
				persistent[k[lenPP:]] = v
			}
		}
	}

	// return original ctx if no invalid key in map
	// make new node, and transfer map to list
	return withNodeFromMaps(ctx, persistent)
}

func newCtxFromMap(ctx context.Context, m map[string]string) context.Context {
	// make new node
	nd := &node{
		persistent: make([]kv, 0, len(m)),
	}

	// insert new kvs from m to node
	for k, v := range m {
		if len(k) == 0 || len(v) == 0 {
			continue
		}
		switch {
		case strings.HasPrefix(k, PrefixPersistent):
			if len(k) > lenPP {
				nd.persistent = append(nd.persistent, kv{key: k[lenPP:], val: v})
			}
		}
	}

	return withNode(ctx, nd)
}

// SaveMetaInfoToMap set key-value pairs from ctx to m while filtering out transient-upstream data.
func SaveMetaInfoToMap(ctx context.Context, m map[string]string) {
	if len(m) == 0 {
		return
	}

	ctx = TransferForward(ctx)
	if n := getNode(ctx); n != nil {
		for _, kv := range n.persistent {
			m[PrefixPersistent+kv.key] = kv.val
		}
	}
}

// newkvtostore new kvstore and converts a kv slice to map.
func newkvtostore(slice []kv) kvstore {
	kvs := newkvstore(len(slice))
	for _, kv := range slice {
		if kv.key != "" && kv.val != "" {
			kvs[kv.key] = kv.val
		}
	}
	return kvs
}
