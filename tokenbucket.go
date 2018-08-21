package core

import "sync"

type TokenBucket struct {
	tokens chan struct{}
	lock   sync.Mutex
}

func NewTokenBucket(numberOfTokens int) *TokenBucket {
	tb := &TokenBucket{
		tokens: make(chan struct{}, numberOfTokens),
		lock:   sync.Mutex{},
	}

	for i := 0; i < numberOfTokens; i++ {
		tb.tokens <- struct{}{}
	}

	return tb
}

func (tb *TokenBucket) TakeToken() {
	tb.lock.Lock()
	defer tb.lock.Unlock()

	<-tb.tokens
}

func (tb *TokenBucket) ReleaseToken() {
	tb.lock.Lock()
	defer tb.lock.Unlock()

	tb.tokens <- struct{}{}
}
