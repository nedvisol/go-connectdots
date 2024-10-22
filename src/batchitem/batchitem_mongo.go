package batchitem

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type batchItemRepository struct {
	collection *mongo.Collection
}

// NewBatchItemRepository creates a new BatchItemRepository.
func NewBatchItemRepository(db *mongo.Database) BatchItemRepository {
	return &batchItemRepository{
		collection: db.Collection("batch_items"), // MongoDB collection name
	}
}

// Create inserts a new BatchItem into the collection.
func (r *batchItemRepository) Create(ctx context.Context, item *BatchItem) error {
	// item.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	// item.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

// FindByID retrieves a BatchItem by its ID.
func (r *batchItemRepository) FindByID(ctx context.Context, id string) (*BatchItem, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var item BatchItem
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &item, nil
}

// Update modifies an existing BatchItem.
func (r *batchItemRepository) Update(ctx context.Context, item *BatchItem) error {
	// item.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	filter := bson.M{"_id": item.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, item)
	return err
}

// Delete removes a BatchItem from the collection.
func (r *batchItemRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// FindAll retrieves all BatchItems.
func (r *batchItemRepository) FindAll(ctx context.Context) ([]*BatchItem, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*BatchItem
	for cursor.Next(ctx) {
		var item BatchItem
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, cursor.Err()
}

// ****** usage

// func main() {
// 	// Set client options
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Replace with your MongoDB URI

// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Check the connection
// 	err = client.Ping(context.TODO(), nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Connected to MongoDB!")

// 	// Use the database and repository
// 	db := client.Database("your_database_name") // Replace with your database name
// 	batchItemRepo := repositories.NewBatchItemRepository(db)

// 	// Example: Create a new BatchItem
// 	item := &models.BatchItem{
// 		Name:     "Item1",
// 		Quantity: 100,
// 	}
// 	err = batchItemRepo.Create(context.TODO(), item)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("BatchItem created:", item)
// }
