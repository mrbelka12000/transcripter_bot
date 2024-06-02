package repo

import (
	"context"
	"transcripter_bot/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	db             *mongo.Database
	collectionName string
}

func NewRepo(db *mongo.Database, collectionName string) *Repo {
	return &Repo{
		db:             db,
		collectionName: collectionName,
	}
}

func (r *Repo) Create(ctx context.Context, item models.Item) (*mongo.InsertOneResult, error) {
	collection := r.db.Collection(r.collectionName)
	return collection.InsertOne(ctx, item)
}

func (r *Repo) Read(ctx context.Context, id string) (*models.Item, error) {
	collection := r.db.Collection(r.collectionName)
	filter := bson.M{"_id": id}
	var item models.Item
	err := collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repo) Search(ctx context.Context, filter bson.M) ([]*models.Item, error) {
	collection := r.db.Collection(r.collectionName)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*models.Item
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repo) Update(ctx context.Context, id string, updatedItem models.Item) (*mongo.UpdateResult, error) {
	collection := r.db.Collection(r.collectionName)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": updatedItem,
	}
	return collection.UpdateOne(ctx, filter, update)
}

func (r *Repo) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	collection := r.db.Collection(r.collectionName)
	filter := bson.M{"_id": id}
	return collection.DeleteOne(ctx, filter)
}
