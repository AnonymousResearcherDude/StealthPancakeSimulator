package types

import (
	"StealthPancakeSimulator/config"
	"sync"
	"testing"
)

const path = "../../../"

func TestAddToCache_Unlimited(t *testing.T) {
	config.InitConfigWithPath(path)
	config.SetCacheModel(0)

	cache := CacheStruct{
		Size:           3,
		CacheMap:       make(CacheMap),
		CacheFreqMap:   make(CacheFreqMap),
		CacheMutex:     &sync.Mutex{},
		EvictionPolicy: GetCachePolicy(),
	}

	for i := 0; i < 100; i++ {
		cache.AddToCache(ChunkId(i), NodeId(1))
	}

	for i := 0; i < 100; i++ {
		if _, exists := cache.CacheMap[ChunkId(i)]; !exists {
			t.Errorf("Expected ChunkId %d to be added", i)
		}
	}
}

func TestAddToCache_Proximity(t *testing.T) {
	config.InitConfig()
	config.SetCacheModel(1)

	cache := CacheStruct{
		Size:           3,
		CacheMap:       make(CacheMap),
		CacheFreqMap:   make(CacheFreqMap),
		CacheMutex:     &sync.Mutex{},
		EvictionPolicy: GetCachePolicy(),
	}

	// 79 = 64 + 8 + 4 + 2 + 1,
	cache.AddToCache(ChunkId(80), NodeId(79)) // 80 = 64 + 16,		dist = 21 (26 - 5)
	cache.AddToCache(ChunkId(76), NodeId(79)) // 76 = 64 + 8 + 4,	dist = 24 (26 - 2)
	cache.AddToCache(ChunkId(00), NodeId(79)) // 0,				dist = 19 (26 - 7)
	cache.AddToCache(ChunkId(64), NodeId(79)) // 64,				dist = 22 (26 - 4)

	if _, exists := cache.CacheMap[ChunkId(76)]; exists {
		t.Errorf("Expected ChunkId 76 to be removed")
	}

	if _, exists := cache.CacheMap[ChunkId(80)]; !exists {
		t.Errorf("Expected ChunkId 80 to be added")
	}
	if _, exists := cache.CacheMap[ChunkId(00)]; !exists {
		t.Errorf("Expected ChunkId 00 to be added")
	}
	if _, exists := cache.CacheMap[ChunkId(64)]; !exists {
		t.Errorf("Expected ChunkId 64 to be added")
	}
}

func TestAddToCache_LeastRecentUsed(t *testing.T) {
	config.InitConfig()
	config.SetCacheModel(2)

	cache := CacheStruct{
		Size:           3,
		CacheMap:       make(CacheMap),
		CacheFreqMap:   make(CacheFreqMap),
		CacheMutex:     &sync.Mutex{},
		EvictionPolicy: GetCachePolicy(),
	}

	cache.AddToCache(ChunkId(1), NodeId(1))
	cache.AddToCache(ChunkId(2), NodeId(1))
	cache.AddToCache(ChunkId(1), NodeId(1))
	cache.AddToCache(ChunkId(3), NodeId(1))
	cache.AddToCache(ChunkId(4), NodeId(1))

	cache.EvictionPolicy.UpdateCacheMap(&cache, ChunkId(4), 0)

	if _, exists := cache.CacheMap[ChunkId(2)]; exists {
		t.Errorf("Expected ChunkId 2 to be removed")
	}

	if _, exists := cache.CacheMap[ChunkId(1)]; !exists {
		t.Errorf("Expected ChunkId 1 to be added")
	}
	if _, exists := cache.CacheMap[ChunkId(3)]; !exists {
		t.Errorf("Expected ChunkId 3 to be added")
	}
	if _, exists := cache.CacheMap[ChunkId(4)]; !exists {
		t.Errorf("Expected ChunkId 4 to be added")
	}
}

func TestAddToCache_LeastFrequentlyUsed(t *testing.T) {
	config.InitConfig()
	config.SetCacheModel(3)

	cache := CacheStruct{
		Size:           3,
		CacheMap:       make(CacheMap),
		CacheFreqMap:   make(CacheFreqMap),
		CacheMutex:     &sync.Mutex{},
		EvictionPolicy: GetCachePolicy(),
	}

	for i := 0; i < 1; i++ {
		cache.AddToCache(ChunkId(1), NodeId(1))
	}
	for i := 0; i < 1; i++ {
		cache.AddToCache(ChunkId(2), NodeId(1))
	}
	for i := 0; i < 1; i++ {
		cache.AddToCache(ChunkId(8), NodeId(1))
	}

	cache.AddToCache(ChunkId(4), NodeId(1))
	cache.AddToCache(ChunkId(1), NodeId(1))

	// 4 should be removed, but there's not way for it to make it into the cache
	if _, exists := cache.CacheMap[ChunkId(2)]; exists {
		t.Errorf("Expected ChunkId 2 to be removed")
	}

	if _, exists := cache.CacheMap[ChunkId(1)]; !exists {
		t.Errorf("Expected ChunkId 1 to be added")
	}
	if _, exists := cache.CacheMap[ChunkId(8)]; !exists {
		t.Errorf("Expected ChunkId 8 to be added")
	}
	if _, exists := cache.CacheMap[ChunkId(4)]; !exists {
		t.Errorf("Expected ChunkId 4 to be added")
	}
}

func TestCacheStruct_Contains(t *testing.T) {
	config.InitConfig()
	config.SetCacheModel(0)

	cache := CacheStruct{
		Size:           3,
		CacheMap:       make(CacheMap),
		CacheFreqMap:   make(CacheFreqMap),
		CacheMutex:     &sync.Mutex{},
		EvictionPolicy: GetCachePolicy(),
	}

	network := &Network{}
	network.Bits = 1
	node := network.node(NodeId(1))

	cache.AddToCache(ChunkId(1), node.Id)
	cache.AddToCache(ChunkId(2), node.Id)
	cache.AddToCache(ChunkId(3), node.Id)

	if !cache.Contains(ChunkId(2)) {
		t.Errorf("Expected ChunkId 2 to be in the cache, but it's not")
	}

	if cache.Contains(ChunkId(4)) {
		t.Errorf("Expected ChunkId 4 to not be in the cache, but it is")
	}
}
