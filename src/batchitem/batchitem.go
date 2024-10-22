package batchitem

import "context"

type BatchItem struct {
	ID        int
	Timestamp int
	URI       string
	//http option?
}

type BatchItemRepository interface {
	Create(ctx context.Context, item *BatchItem) error
	FindByID(ctx context.Context, id string) (*BatchItem, error)
	Update(ctx context.Context, item *BatchItem) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*BatchItem, error)
}
