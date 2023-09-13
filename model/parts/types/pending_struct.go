package types

import (
	"StealthPancakeSimulator/config"
	"sync"
)

type QueuedChunk struct {
	ChunkId   ChunkId
	Counter   int
	LastEpoch int
}

type PendingStruct struct {
	PendingQueue []QueuedChunk
	CurrentIndex int
	PendingMutex *sync.Mutex
}

func (p *PendingStruct) AddPendingChunkId(chunkId ChunkId, curEpoch int) bool {
	p.PendingMutex.Lock()
	defer p.PendingMutex.Unlock()
	chunkIdIndex := p.GetQueuedChunkIndex(chunkId)
	isNewChunk := false
	if chunkIdIndex == -1 { // new chunk
		newChunkStruct := QueuedChunk{
			ChunkId:   chunkId,
			Counter:   0,
			LastEpoch: curEpoch,
		}
		p.PendingQueue = append([]QueuedChunk{newChunkStruct}, p.PendingQueue...)
		isNewChunk = true

	} else { // chunk seen before
		if p.PendingQueue[chunkIdIndex].Counter < config.GetBinSize() {
			p.PendingQueue[chunkIdIndex].Counter++

		} else { // remove queued chunk
			p.PendingQueue = append(p.PendingQueue[:chunkIdIndex], p.PendingQueue[chunkIdIndex+1:]...)
		}
	}
	return isNewChunk
}

func (p *PendingStruct) DeletePendingChunkId(chunkId ChunkId) {
	p.PendingMutex.Lock()
	defer p.PendingMutex.Unlock()

	chunkIdIndex := p.GetQueuedChunkIndex(chunkId)
	if chunkIdIndex != -1 {
		p.PendingQueue = append(p.PendingQueue[:chunkIdIndex], p.PendingQueue[chunkIdIndex+1:]...)
	}

}

func (p *PendingStruct) GetChunkFromQueue(curEpoch int) (QueuedChunk, bool) {
	p.PendingMutex.Lock()
	defer p.PendingMutex.Unlock()

	for i := 0; i < len(p.PendingQueue); i++ {
		currentIndex := p.GetAndUpdateCurrentIndex()
		chunkFrontOfQueue := p.PendingQueue[currentIndex]

		if chunkFrontOfQueue.LastEpoch < curEpoch {
			p.PendingQueue[currentIndex].LastEpoch = curEpoch

			return chunkFrontOfQueue, true
		}
	}

	return QueuedChunk{}, false
}

func (p *PendingStruct) GetQueuedChunkIndex(chunkId ChunkId) int {

	for i, v := range p.PendingQueue {
		if v.ChunkId == chunkId {
			return i
		}
	}
	return -1
}

func (p *PendingStruct) GetAndUpdateCurrentIndex() int {

	p.CurrentIndex--
	if p.CurrentIndex < 0 {
		p.CurrentIndex = len(p.PendingQueue) - 1
	}
	return p.CurrentIndex
}

//func (p *PendingStruct) UpdateEpoch(chunkId ChunkId, curEpoch int) int {
//	p.PendingMutex.Lock()
//	defer p.PendingMutex.Unlock()
//
//	chunkStructIndex := p.GetQueuedChunkIndex(chunkId)
//	if chunkStructIndex != -1 {
//		p.PendingQueue[chunkStructIndex].LastEpoch = curEpoch
//		return p.PendingQueue[chunkStructIndex].LastEpoch
//	}
//
//	return -1
//}

//type PendingMap map[NodeId]PendingQueue

//type PendingStruct struct {
//	PendingMap           PendingMap
//	PendingMutex         *sync.Mutex
//	UniquePendingCounter int32
//}

//func (p *PendingStruct) GetPending(originator NodeId) (PendingQueue, bool) {
//	p.PendingMutex.Lock()
//	defer p.PendingMutex.Unlock()
//	pending, ok := p.PendingMap[originator]
//	if ok {
//		return pending, true
//	}
//	return PendingQueue{ChunkQueue: []QueuedChunk{}, CurrentIndex: 0}, false
//}

//func (p *PendingStruct) AddPendingChunkId(originator NodeId, chunkId ChunkId, curEpoch int) {
//	pending, _ := p.GetPending(originator)
//	chunkStructIndex := p.GetChunkStructIndex(pending.ChunkQueue, chunkId)
//
//	p.PendingMutex.Lock()
//	defer p.PendingMutex.Unlock()
//
//	if chunkStructIndex == -1 { // new chunk
//		newChunkStruct := QueuedChunk{
//			ChunkId:   chunkId,
//			Counter:   0,
//			LastEpoch: curEpoch,
//		}
//		pending.ChunkQueue = append([]QueuedChunk{newChunkStruct}, pending.ChunkQueue...)
//		p.PendingMap[originator] = pending
//		p.UniquePendingCounter++
//
//	} else { // chunk seen before
//		if pending.ChunkQueue[chunkStructIndex].Counter < constants.GetBinSize() {
//			pending.ChunkQueue[chunkStructIndex].Counter++
//			p.PendingMap[originator] = pending
//
//		} else { // remove chunkStruct
//			pending.ChunkQueue = append(pending.ChunkQueue[:chunkStructIndex], pending.ChunkQueue[chunkStructIndex+1:]...)
//			if len(pending.ChunkQueue) == 0 {
//				delete(p.PendingMap, originator)
//			}
//		}
//	}
//}

//func (p *PendingStruct) DeletePendingChunkId(originator NodeId, chunkId ChunkId) {
//	pending, _ := p.GetPending(originator)
//
//	if len(pending.ChunkQueue) > 0 {
//		chunkIdIndex := p.GetChunkStructIndex(pending.ChunkQueue, chunkId)
//		if chunkIdIndex != -1 {
//			p.PendingMutex.Lock()
//			defer p.PendingMutex.Unlock()
//			pending.ChunkQueue = append(pending.ChunkQueue[:chunkIdIndex]) // Delete chunk front of queue
//			if len(pending.ChunkQueue) == 0 {
//				delete(p.PendingMap, originator)
//				return
//			} else {
//				p.PendingMap[originator] = pending
//			}
//		}
//	}
//}

//func (p *PendingStruct) GetChunkFromQueue(originator NodeId) (QueuedChunk, bool) {
//	pending, ok := p.GetPending(originator)
//	if ok {
//		p.PendingMutex.Lock()
//		defer p.PendingMutex.Unlock()
//		currentIndex := p.GetAndUpdateCurrentIndex(pending, originator)
//		chunkFrontOfQueue := pending.ChunkQueue[currentIndex]
//		return chunkFrontOfQueue, true
//	}
//	return QueuedChunk{}, false
//}

//func (p *PendingStruct) UpdateEpoch(originator NodeId, chunkId ChunkId, curEpoch int) int {
//	p.PendingMutex.Lock()
//	defer p.PendingMutex.Unlock()
//	pending, ok := p.PendingMap[originator]
//	if ok {
//		chunkStructIndex := p.GetChunkStructIndex(pending.ChunkQueue, chunkId)
//		if chunkStructIndex != -1 {
//			p.PendingMap[originator].ChunkQueue[chunkStructIndex].LastEpoch = curEpoch
//			return p.PendingMap[originator].ChunkQueue[chunkStructIndex].LastEpoch
//		}
//	}
//	return -1
//}

//func (p *PendingStruct) GetAndUpdateCurrentIndex(pending PendingQueue, originator NodeId) int {
//
//	pending.CurrentIndex--
//	if pending.CurrentIndex < 0 || pending.CurrentIndex >= len(pending.ChunkQueue) {
//		pending.CurrentIndex = len(pending.ChunkQueue) - 1
//		if pending.CurrentIndex < 0 {
//			pending.CurrentIndex = 0
//		}
//	}
//	p.PendingMap[originator] = pending
//	return pending.CurrentIndex
//}

//func (p *PendingStruct) SetEpochDecrement(originator int) int32 {
//	p.PendingMutex.Lock()
//	defer p.PendingMutex.Unlock()
//	pending, ok := p.PendingMap[originator]
//	if ok {
//		pending.EpokeDecrement = int32(len(pending.ChunkQueue))
//		p.PendingMap[originator] = pending
//		return pending.EpokeDecrement
//	}
//	return -1
//}

//func (p *PendingStruct) AddPending(originator int, chunkId int) {
//	p.PendingMutex.Lock()
//	pendingNode := p.PendingMap[originator]
//	pendingNode.ChunkIds = append(pendingNode.ChunkIds, chunkId)
//	pendingNode.PendingCounter = 1
//	p.PendingMap[originator] = pendingNode
//	p.PendingMutex.Unlock()
//}
//
//func (p *PendingStruct) AddToPendingQueue(originator int, chunkId int) {
//	p.PendingMutex.Lock()
//	pendingNode := p.PendingMap[originator]
//	pendingNode.ChunkIds = append(pendingNode.ChunkIds, chunkId)
//	pendingNode.PendingCounter++
//	p.PendingMap[originator] = pendingNode
//	p.PendingMutex.Unlock()
//}
//func (p *PendingStruct) DeletePending(originator int) {
//	p.PendingMutex.Lock()
//	delete(p.PendingMap, originator)
//	p.PendingMutex.Unlock()
//}

//func (p *PendingStruct) IncrementPending(originator int) {
//	p.PendingMutex.Lock()
//	pendingNode := p.PendingMap[originator]
//	pendingNode.PendingCounter++
//	p.PendingMap[originator] = pendingNode
//	p.PendingMutex.Unlock()
//}

//func (p *PendingStruct) IsEmpty(originator int) bool {
//	pending := p.GetPending(originator)
//	if len(pending.ChunkIds) > 0 {
//		return false
//	}
//	return true
//}
//
//func (p *PendingStruct) DeletePendingNodeId(originator int, pendingNodeIdIndex int) {
//	p.PendingMutex.Lock()
//	pendingNode := p.PendingMap[originator]
//	pendingNode.ChunkIds = append(pendingNode.ChunkIds[:pendingNodeIdIndex])
//	p.PendingMap[originator] = pendingNode
//	p.PendingMutex.Unlock()
//}
