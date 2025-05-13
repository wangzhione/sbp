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

func TimeoutLock(key string, timeout time.Duration) bool {
	return defaultLocalLock.mutex(key).TimeoutLock(timeout)
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
	// 如果不能 断言 会 panic
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

func (l *LocalLock) TimeoutLock(key string, timeout time.Duration) bool {
	return l.mutex(key).TimeoutLock(timeout)
}

// Unlock 解锁
func (l *LocalLock) Unlock(key string) {
	if val, ok := l.local.Load(key); ok {
		val.(*Locker).Unlock()
	} else {
		// 提醒 case
		println("panic: multiple Unlock")
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

func (k *Locker) TimeoutLock(timeout time.Duration) bool {
	select {
	case <-k.ch:
		return true
	default:
	}

	timer := time.NewTimer(timeout) // 请用新一点 Go 版本, timer 相关操作更安全, 老的不再维护, 需要可翻阅当前文件老代码
	defer timer.Stop()

	select {
	case <-k.ch:
		return true
	case <-timer.C:
		return false
	}
}
