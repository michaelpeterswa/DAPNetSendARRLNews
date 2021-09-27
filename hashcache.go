package main

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

type HashCache struct {
	ttl float64

	mu     sync.Mutex
	hashes map[string]*time.Time
}

type Hash struct {
	Hash string
	Time *time.Time
}

func (hc *HashCache) Init(ttl float64) {
	hc.hashes = make(map[string]*time.Time)
	hc.ttl = ttl
	logger.Info("Initialized HashCache...")
}

func (hc *HashCache) Set(hash Hash) {
	hc.mu.Lock()
	hc.hashes[hash.Hash] = hash.Time
	hc.mu.Unlock()
	logger.Info("Created Entry into HashCache", zap.String("entry", hash.Hash))
}

func (hc *HashCache) Get(key string) Hash {
	var hash Hash
	hc.mu.Lock()
	defer hc.mu.Unlock()

	val, ok := hc.hashes[key]
	if ok {
		hash.Hash = key
		hash.Time = val
	}
	return hash
}

func (hc *HashCache) Exists(key string) bool {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	_, ok := hc.hashes[key]
	return ok
}

func (hc *HashCache) Clean() {
	cnt := 0
	hc.mu.Lock()
	for key, value := range hc.hashes {
		if time.Since(*value).Hours() > hc.ttl {
			delete(hc.hashes, key)
			cnt++
		}
	}
	hc.mu.Unlock()
	logger.Info("Cache Cleaned", zap.Int("entriesRemoved", cnt))
}

func (hc *HashCache) Pprint() {
	hc.mu.Lock()
	for key, value := range hc.hashes {
		fmt.Printf("hash: %s || date: %s\n", key, value)
	}
	hc.mu.Unlock()
}
