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
	// if m.isLocked {
	// log.Printf("[System] [Mutex] Tried locking a locked mutex")
	// return
	// }

	// log.Printf("[System] [Mutex] Locking mutex")
	m.mutex.Lock()
	m.isLocked = true
}

func (m *Mutex) Unlock(caller string) {
	if !m.isLocked {
		log.Printf("[System] [Mutex] Tried unlocking a unlocked mutex by %v", caller)
		return
	}

	m.isLocked = false
	// log.Printf("[System] [Mutex] Unlocking mutex")
	m.mutex.Unlock()
}
