package storage

// ReceiptStorage defines the interface for receipt storage
type ReceiptStorage interface {
	SaveReceipt(id string, points int) error
	GetPoints(id string) (int, bool)
}
