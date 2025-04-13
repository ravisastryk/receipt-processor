package storage

import (
	"sync"
)

// MemoryStorage implements ReceiptStorage using an in-memory map
type MemoryStorage struct {
	receipts map[string]int
	mutex    sync.RWMutex // TODO we can have separate mutex for Read vs Write for perf optimization
}

// NewMemoryStorage creates a new memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		receipts: make(map[string]int),
	}
}

// SaveReceipt saves a receipt ID and its points
func (m *MemoryStorage) SaveReceipt(id string, points int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.receipts[id] = points
	return nil
}

// GetPoints retrieves the points for a receipt
func (m *MemoryStorage) GetPoints(id string) (int, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	points, exists := m.receipts[id]
	return points, exists
}
