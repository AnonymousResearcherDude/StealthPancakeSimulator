package types

import (
	"StealthPancakeSimulator/config"
	"math/big"
	"sync"

	"github.com/seehuhn/mt19937"
)

type CidStruct struct {
	Cid            string
	ReturenedBytes int
}

type CidQueueStruct struct {
	CidQueue      []CidStruct
	ChunksFromCid []ChunkId
	CidQueueMutex *sync.Mutex
}

func (c *CidQueueStruct) AddToCidQueue(cidStruct CidStruct) {
	c.CidQueueMutex.Lock()
	defer c.CidQueueMutex.Unlock()

	c.CidQueue = append(c.CidQueue, cidStruct)

	c.FillChunksFromCid()
}

func (c *CidQueueStruct) FillChunksFromCid() {
	if len(c.ChunksFromCid) < 1 && len(c.CidQueue) > 0 {
		firstCidItem := c.FirstCidItem()
		if firstCidItem == nil {
			return
		}
		numberOfChunks := GetNumberOfChunksFromBytes(firstCidItem.ReturenedBytes)
		chunksFromCid := GenerateRandomChunkIdsFromCid(firstCidItem.Cid, numberOfChunks)
		c.ChunksFromCid = chunksFromCid
	}
}

func (c *CidQueueStruct) GetChunkFromCidQueue() ChunkId {
	c.FillChunksFromCid()
	if len(c.ChunksFromCid) > 0 {
		chunkId := c.ChunksFromCid[0]
		c.ChunksFromCid = c.ChunksFromCid[1:]
		return chunkId
	}

	return -1
}

func (c *CidQueueStruct) FirstCidItem() *CidStruct {
	if len(c.CidQueue) > 0 {
		firstCid := &c.CidQueue[0]
		c.CidQueue = c.CidQueue[1:]
		return firstCid
	}

	return nil
}

func GetNumberOfChunksFromBytes(bytesReturned int) int {
	kBytes := int(bytesReturned/1024.0) + 1
	chunks := (kBytes / 4) + 1
	return chunks
}

func GenerateRandomChunkIdsFromCid(cid string, numberOfChunks int) []ChunkId {
	allChunks := make([]ChunkId, 0, numberOfChunks)

	seedInt := new(big.Int).SetBytes([]byte(cid))
	source := mt19937.New()
	source.Seed(seedInt.Int64())

	max := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(config.GetBits())), nil)
	for i := 0; i < numberOfChunks; i++ {
		num := new(big.Int)
		num.SetUint64(source.Uint64())
		number := num.Mod(num, max)
		allChunks = append(allChunks, ChunkId(number.Int64()))
	}

	return allChunks
}
