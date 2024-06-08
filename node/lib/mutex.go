package node

import (
	"log"
	"sync"
)

type Mutex struct {
	isLocked bool
	mutex    sync.Mutex
}

func (m *Mutex) Lock() {
	if m.isLocked {
		log.Printf("[System] [Mutex] Tried locking a locked mutex")
		return
	}

	m.isLocked = true
	// log.Printf("[System] [Mutex] Locking mutex")
	m.mutex.Lock()
}

func (m *Mutex) Unlock() {
	if !m.isLocked {
		log.Printf("[System] [Mutex] Tried unlocking a unlocked mutex")
		return
	}

	m.isLocked = false
	// log.Printf("[System] [Mutex] Unlocking mutex")
	m.mutex.Unlock()
}