package metainfo

import "sync"

type kvstore map[string]string

var kvstorep sync.Pool

func newkvstore(size ...int) kvstore {
	kvs := kvstorep.Get()
	if kvs != nil {
		return kvs.(kvstore)
	}

	if len(size) > 0 && size[0] > 0 {
		return make(kvstore, size[0])
	}
	return make(kvstore)
}

func (store kvstore) size() int {
	return len(store)
}

func (store kvstore) recycle() {
	clear(store)
	kvstorep.Put(store)
}
