package cacheditem

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type cachedItemMongoRepository struct {
	collection *mongo.Collection
}

// NewCachedItemRepository creates a new CachedItem repository instance.
func NewCachedItemMongoRepository(db *mongo.Database) CachedItemRepository {
	return &cachedItemMongoRepository{
		collection: db.Collection("cached_items"), // MongoDB collection name
	}
}

// Create inserts a new CachedItem into the collection.
func (r *cachedItemMongoRepository) Create(ctx context.Context, item *CachedItem) error {
	item.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

// FindByKey finds a CachedItem by its cache key.
func (r *cachedItemMongoRepository) FindByKey(ctx context.Context, key string) (*CachedItem, error) {
	var item CachedItem
	err := r.collection.FindOne(ctx, bson.M{"key": key}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No item found
		}
		return nil, err
	}
	return &item, nil
}

// Update modifies an existing CachedItem in the collection.
func (r *cachedItemMongoRepository) Update(ctx context.Context, item *CachedItem) error {
	filter := bson.M{"_id": item.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, item)
	return err
}

// DeleteByKey removes a CachedItem from the collection by its key.
func (r *cachedItemMongoRepository) DeleteByKey(ctx context.Context, key string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"key": key})
	return err
}

// FindAll retrieves all CachedItems from the collection.
func (r *cachedItemMongoRepository) FindAll(ctx context.Context) ([]*CachedItem, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*CachedItem
	for cursor.Next(ctx) {
		var item CachedItem
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, cursor.Err()
}
