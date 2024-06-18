package storage

import "context"

type ListReturnFunc func(any) string

type Storer interface {
	Save(batchDir, filename string, data []string) error
	Read(key string) ([]string, error)
	ReadBatch(readItems []ReadItem) (map[string][]string, error)
	Delete(key string) error
	List(ctx context.Context, delim, blockBatch string, timeout int, returnFunc ListReturnFunc) ([]string, error)
}

type ReadItem struct {
	Key    string
	RowIds []uint64
}
