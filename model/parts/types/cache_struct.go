package types

import (
	"go-incentive-simulation/config"
	"go-incentive-simulation/model/general"
	"sync"
	"time"

	"github.com/zavitax/sortedset-go"
)

type CacheMap map[ChunkId]CacheData
type CacheFreqMap map[ChunkId]int

type CacheData struct {
	Proximity    int
	LastTimeUsed time.Time
	Frequency    int
}

type CacheStruct struct {
	Size           int
	Node           *Node
	CacheMap       CacheMap
	CacheFreqMap   CacheFreqMap
	CacheMutex     *sync.Mutex
	EvictionPolicy CachePolicy
}

type CachePolicy interface {
	UpdateCacheMap(c *CacheStruct, newChunkId ChunkId, distance int)
}

type (
	unlimitedPolicy struct {
	}

	proximityPolicy struct {
		ChunkSet *sortedset.SortedSet[ChunkId, int, CacheData]
	}

	lruPolicy struct {
		ChunkSet *sortedset.SortedSet[ChunkId, int, CacheData]
	}

	lfuPolicy struct {
		ChunkSet *sortedset.SortedSet[ChunkId, int, CacheData]
	}
)

func GetCachePolicy() CachePolicy {
	policy := config.GetCacheModel()
	if policy == -1 {
		return nil
	}

	if policy == 0 {
		return &unlimitedPolicy{}
	} else if policy == 1 {
		return &proximityPolicy{ChunkSet: sortedset.New[ChunkId, int, CacheData]()}
	} else if policy == 2 {
		return &lruPolicy{ChunkSet: sortedset.New[ChunkId, int, CacheData]()}
	} else if policy == 3 {
		return &lfuPolicy{ChunkSet: sortedset.New[ChunkId, int, CacheData]()}
	} else {
		return nil
	}
}

func FindDistance(chunkId ChunkId, nodeId NodeId) int {
	return config.GetBits() - general.BitLength(chunkId.ToInt()^nodeId.ToInt())
}

func (c *CacheStruct) AddToCache(chunkId ChunkId, nodeId NodeId) CacheMap {
	c.CacheMutex.Lock()
	defer c.CacheMutex.Unlock()
	distance := FindDistance(chunkId, nodeId)

	freq := 1
	if _, ok := c.CacheFreqMap[chunkId]; ok {
		c.CacheFreqMap[chunkId]++
		freq = c.CacheFreqMap[chunkId]
	} else {
		c.CacheFreqMap[chunkId] = 1
	}

	if _, ok := c.CacheMap[chunkId]; ok {
		currData := c.CacheMap[chunkId]
		currData.Frequency++
		currData.Proximity = distance
		c.CacheMap[chunkId] = currData
	} else {
		newCacheData := CacheData{
			Proximity:    distance,
			LastTimeUsed: time.Now(),
			Frequency:    freq,
		}
		c.CacheMap[chunkId] = newCacheData
	}

	if c.EvictionPolicy != nil {
		c.EvictionPolicy.UpdateCacheMap(c, chunkId, distance)
	}

	return c.CacheMap
}

func (p *unlimitedPolicy) UpdateCacheMap(c *CacheStruct, newChunkId ChunkId, distance int) {
}

func (p *proximityPolicy) UpdateCacheMap(c *CacheStruct, newChunkId ChunkId, distance int) {
	p.ChunkSet.AddOrUpdate(newChunkId, distance, CacheData{distance, time.Now(), 1})

	if len(c.CacheMap) <= int(c.Size) {
		return
	}

	chunkIdToDelete := p.ChunkSet.PopMax().Key()
	delete(c.CacheMap, chunkIdToDelete)
}

func (p *lruPolicy) UpdateCacheMap(c *CacheStruct, newChunkId ChunkId, distance int) {
	p.ChunkSet.AddOrUpdate(newChunkId, time.Now().Nanosecond(), CacheData{distance, time.Now(), 1})

	if len(c.CacheMap) <= int(c.Size) {
		return
	}

	chunkIdToDelete := p.ChunkSet.PopMin().Key()
	delete(c.CacheMap, chunkIdToDelete)
}

func (p *lfuPolicy) UpdateCacheMap(c *CacheStruct, newChunkId ChunkId, distance int) {
	prev := p.ChunkSet.GetByKey(newChunkId)
	freq := 1
	if prev != nil {
		freq = prev.Value.Frequency + 1
	}

	if _, exists := c.CacheFreqMap[newChunkId]; exists {
		freq = c.CacheFreqMap[newChunkId]
	}

	p.ChunkSet.AddOrUpdate(newChunkId, freq, CacheData{distance, time.Now(), freq})

	if len(c.CacheMap) <= int(c.Size) {
		return
	}

	itemToDelete := p.ChunkSet.PopMin()
	if itemToDelete.Key() == newChunkId {
		itemToDelete, _ = p.ChunkSet.PopMin(), p.ChunkSet.AddOrUpdate(itemToDelete.Key(), itemToDelete.Score(), itemToDelete.Value)
	}
	delete(c.CacheMap, itemToDelete.Key())
}

func (c *CacheStruct) Contains(chunkId ChunkId) bool {
	c.CacheMutex.Lock()
	defer c.CacheMutex.Unlock()

	if _, ok := c.CacheMap[chunkId]; ok {
		cacheData := c.CacheMap[chunkId]
		cacheData.LastTimeUsed = time.Now()
		cacheData.Frequency = cacheData.Frequency + 1
		c.CacheMap[chunkId] = cacheData
		return true
	}

	return false
}
