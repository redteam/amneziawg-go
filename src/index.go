package main

import (
	"crypto/rand"
	"sync"
)

/* Index=0 is reserved for unset indecies
 *
 */

type IndexTable struct {
	mutex      sync.RWMutex
	keypairs   map[uint32]*KeyPair
	handshakes map[uint32]*Handshake
}

func randUint32() (uint32, error) {
	var buff [4]byte
	_, err := rand.Read(buff[:])
	id := uint32(buff[0])
	id <<= 8
	id |= uint32(buff[1])
	id <<= 8
	id |= uint32(buff[2])
	id <<= 8
	id |= uint32(buff[3])
	return id, err
}

func (table *IndexTable) Init() {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	table.keypairs = make(map[uint32]*KeyPair)
	table.handshakes = make(map[uint32]*Handshake)
}

func (table *IndexTable) NewIndex(handshake *Handshake) (uint32, error) {
	table.mutex.Lock()
	defer table.mutex.Unlock()
	for {
		// generate random index

		id, err := randUint32()
		if err != nil {
			return id, err
		}
		if id == 0 {
			continue
		}

		// check if index used

		_, ok := table.keypairs[id]
		if ok {
			continue
		}
		_, ok = table.handshakes[id]
		if ok {
			continue
		}

		// update the index

		delete(table.handshakes, handshake.localIndex)
		handshake.localIndex = id
		table.handshakes[id] = handshake
		return id, nil
	}
}

func (table *IndexTable) LookupKeyPair(id uint32) *KeyPair {
	table.mutex.RLock()
	defer table.mutex.RUnlock()
	return table.keypairs[id]
}

func (table *IndexTable) LookupHandshake(id uint32) *Handshake {
	table.mutex.RLock()
	defer table.mutex.RUnlock()
	return table.handshakes[id]
}