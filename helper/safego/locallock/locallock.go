// Package locallock provides a local locking mechanism for synchronizing access to resources by key.
package locallock

import (
	"time"

	"github.com/wangzhione/sbp/structs/maps"
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

func TimeoutLock(key string, timeout time.Duration) bool {
	return defaultLocalLock.mutex(key).TimeoutLock(timeout)
}

// Unlock 解锁
func Unlock(key string) {
	defaultLocalLock.mutex(key).Unlock()
}

type LocalLock struct {
	local maps.Map[string, *Locker]
}

// 内部方法：获取或创建对应 key 的锁
func (l *LocalLock) mutex(key string) *Locker {
	actual, _ := l.local.LoadOrStore(key, NewLocker())
	return actual
}

// Lock 阻塞式加锁
func (l *LocalLock) Lock(key string) {
	l.mutex(key).Lock()
}

// TryLock 非阻塞尝试加锁
func (l *LocalLock) TryLock(key string) bool {
	return l.mutex(key).TryLock()
}

func (l *LocalLock) TimeoutLock(key string, timeout time.Duration) bool {
	return l.mutex(key).TimeoutLock(timeout)
}

// Unlock 解锁
func (l *LocalLock) Unlock(key string) {
	if actual, ok := l.local.Load(key); ok {
		actual.Unlock()
	} else {
		// 提醒 case
		println("multiple Unlock: " + key)
	}
}

type Locker struct {
	Semaphore chan struct{}
}

func NewLocker() *Locker {
	k := &Locker{
		Semaphore: make(chan struct{}, 1),
	}
	k.Semaphore <- struct{}{} // 初始为可用状态
	return k
}

func (k *Locker) Lock() {
	<-k.Semaphore
}

func (k *Locker) Unlock() {
	select {
	case k.Semaphore <- struct{}{}:
	default:
		// 强制要求 Lock 和 UnLock 一一对应 & 不支持嵌套 Lock + Unlock
		panic("error: Locker unlock of unlocked mutex")
	}
}

func (k *Locker) TryLock() bool {
	select {
	case <-k.Semaphore:
		return true
	default:
		return false
	}
}

func (k *Locker) TimeoutLock(timeout time.Duration) bool {
	select {
	case <-k.Semaphore:
		return true
	default:
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-k.Semaphore:
		return true
	case <-timer.C:
		return false
	}
}
