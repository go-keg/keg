package pubsub

import (
	"sync"
)

type PubSub[T any] struct {
	mu          sync.RWMutex
	subscribers map[int]map[chan T]struct{} // accountID -> channels
	allSubs     map[chan T]int              // chan -> accountID
}

func New[T any]() *PubSub[T] {
	return &PubSub[T]{
		subscribers: make(map[int]map[chan T]struct{}),
		allSubs:     make(map[chan T]int),
	}
}

// Subscribe 用户订阅，返回其消息通道
func (r *PubSub[T]) Subscribe(accountID int) chan T {
	ch := make(chan T, 16)
	r.mu.Lock()
	defer r.mu.Unlock()

	if accountID != 0 {
		if _, ok := r.subscribers[accountID]; !ok {
			r.subscribers[accountID] = make(map[chan T]struct{})
		}
		r.subscribers[accountID][ch] = struct{}{}
	}
	r.allSubs[ch] = accountID

	return ch
}

// Unsubscribe 取消订阅
func (r *PubSub[T]) Unsubscribe(ch chan T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	accountID, ok := r.allSubs[ch]
	if !ok {
		return
	}
	delete(r.allSubs, ch)

	if subs, ok := r.subscribers[accountID]; ok {
		delete(subs, ch)
		if len(subs) == 0 {
			delete(r.subscribers, accountID)
		}
	}
	close(ch)
}

// Publish 向所有订阅者广播
func (r *PubSub[T]) Publish(msg T) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for ch := range r.allSubs {
		r.safeSend(ch, msg)
	}
}

func (r *PubSub[T]) sendTo(accountID int, msg T) {
	if subs, ok := r.subscribers[accountID]; ok {
		for ch := range subs {
			r.safeSend(ch, msg)
		}
	}
}

// SendTo 发送给单个用户
func (r *PubSub[T]) SendTo(accountID int, msg T) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.sendTo(accountID, msg)
}

// SendToMany 发送给多个用户
func (r *PubSub[T]) SendToMany(accountIDs []int, msg T) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, accountID := range accountIDs {
		r.sendTo(accountID, msg)
	}
}

// safeSend 防止 panic 的发送函数
func (r *PubSub[T]) safeSend(ch chan T, msg T) {
	go func() {
		defer func() { _ = recover() }()
		select {
		case ch <- msg:
		default:
		}
	}()
}
