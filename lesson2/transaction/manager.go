package transaction

import (
	"sync"
	"time"
)

type WALEntry []byte

type TransactionManager struct {
	snapshot []byte
	wal      []WALEntry
	queue    chan []byte
	mutex    sync.Mutex
}

func NewTransactionManager(queue chan []byte) *TransactionManager {
	tm := &TransactionManager{
		snapshot: nil,
		wal:      []WALEntry{},
		queue:    queue,
	}

	go tm.LogsReader()

	go tm.updateTM()

	return tm
}

func (tm *TransactionManager) LogsReader() {
	for {
		entry := <-tm.queue

		tm.mutex.Lock()
		tm.wal = append(tm.wal, entry)
		tm.mutex.Unlock()
	}
}

func (tm *TransactionManager) updateTM() {
	for {
		time.Sleep(time.Minute)

		tm.mutex.Lock()
		if len(tm.wal) != 0 {
			tm.snapshot = tm.wal[len(tm.wal)-1]
			tm.wal = tm.wal[:0]
		}
		tm.mutex.Unlock()
	}
}
