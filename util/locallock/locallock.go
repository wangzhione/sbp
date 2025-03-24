package locallock

import (
	"sync"
)

type LocalLock struct {
	local sync.Map // map[string]*sync.Mutex
}

// 内部方法：获取或创建对应 key 的锁
func (l *LocalLock) mutex(key string) *sync.Mutex {
	actual, _ := l.local.LoadOrStore(key, &sync.Mutex{})
	return actual.(*sync.Mutex)
}

// Lock 阻塞式加锁
func (l *LocalLock) Lock(key string) {
	l.mutex(key).Lock()
}

// TryLock 非阻塞尝试加锁
func (l *LocalLock) TryLock(key string) bool {
	return l.mutex(key).TryLock() // Go 1.18+
}

// Unlock 解锁
func (l *LocalLock) Unlock(key string) {
	if val, ok := l.local.Load(key); ok {
		val.(*sync.Mutex).Unlock()
	}
}

var defaultLocalLock LocalLock

// Lock 阻塞式加锁
func Lock(key string) {
	defaultLocalLock.mutex(key).Lock()
}

// TryLock 非阻塞尝试加锁
func TryLock(key string) bool {
	return defaultLocalLock.mutex(key).TryLock()
}

// Unlock 解锁
func Unlock(key string) {
	defaultLocalLock.mutex(key).Unlock()
}
