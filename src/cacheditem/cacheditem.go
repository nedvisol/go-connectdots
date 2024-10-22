package cacheditem

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CachedItem represents an item to be cached.
type CachedItem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"` // MongoDB Object ID
	Key          string             `bson:"key"`           // The cache key
	Value        string             `bson:"value"`         // The cached value
	ExpiresAtSec int64              `bson:"expires_at"`    // Expiration time for the cache item in seconds since epoch
}

type CachedItemRepository interface {
	Create(ctx context.Context, item *CachedItem) error
	FindByKey(ctx context.Context, key string) (*CachedItem, error)
	Update(ctx context.Context, item *CachedItem) error
	DeleteByKey(ctx context.Context, key string) error
	FindAll(ctx context.Context) ([]*CachedItem, error)
}
