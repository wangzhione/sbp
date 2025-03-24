package locallock

import (
	"sync"
	"time"
)

var defaultLocalLock LocalLock

// Lock 阻塞式加锁
func Lock(key string) {
	defaultLocalLock.mutex(key).Lock()
}

// TryLock 非阻塞尝试加锁
func TryLock(key string) bool {
	return defaultLocalLock.mutex(key).TryLock()
}

func TryLockWithTimeout(key string, timeout time.Duration) bool {
	return defaultLocalLock.mutex(key).TryLockWithTimeout(timeout)
}

// Unlock 解锁
func Unlock(key string) {
	defaultLocalLock.mutex(key).Unlock()
}

type LocalLock struct {
	local sync.Map // map[string]*Locker
}

// 内部方法：获取或创建对应 key 的锁
func (l *LocalLock) mutex(key string) *Locker {
	actual, _ := l.local.LoadOrStore(key, NewLocker())
	return actual.(*Locker)
}

// Lock 阻塞式加锁
func (l *LocalLock) Lock(key string) {
	l.mutex(key).Lock()
}

// TryLock 非阻塞尝试加锁
func (l *LocalLock) TryLock(key string) bool {
	return l.mutex(key).TryLock() // Go 1.18+
}

func (l *LocalLock) TryLockWithTimeout(key string, timeout time.Duration) bool {
	return l.mutex(key).TryLockWithTimeout(timeout)
}

// Unlock 解锁
func (l *LocalLock) Unlock(key string) {
	if val, ok := l.local.Load(key); ok {
		val.(*Locker).Unlock()
	}
}

type Locker struct {
	ch chan struct{}
}

func NewLocker() *Locker {
	k := &Locker{
		ch: make(chan struct{}, 1),
	}
	k.ch <- struct{}{} // 初始为可用状态
	return k
}

func (k *Locker) Lock() {
	<-k.ch
}

func (k *Locker) Unlock() {
	select {
	case k.ch <- struct{}{}:
	default:
		// 强制要求 Lock 和 UnLock 一一对应 & 不支持嵌套 Lock + Unlock
		panic("error: Locker unlock of unlocked mutex")
	}
}

func (k *Locker) TryLock() bool {
	select {
	case <-k.ch:
		return true
	default:
		return false
	}
}

func (k *Locker) TryLockWithTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-k.ch:
		return true
	case <-timer.C:
		return false
	}
}
