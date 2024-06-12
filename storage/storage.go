package storage

type Storer interface {
	Save(batchDir, filename string, data []string) error
	Read(key string) ([]string, error)
	ReadBatch(readItems []ReadItem) (map[string][]string, error)
	Delete(key string) error
}

type ReadItem struct {
	Key    string
	RowIds []uint64
}
