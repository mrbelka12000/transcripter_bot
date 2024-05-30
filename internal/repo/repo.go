package repo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Item struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
}

func Create(db *mongo.Database, item Item) (*mongo.InsertOneResult, error) {
	collection := db.Collection("items")
	insertResult, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	return insertResult, nil
}

func Read(db *mongo.Database, id string) (Item, error) {
	collection := db.Collection("items")
	filter := bson.M{"_id": id}
	var result Item
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Print(err.Error())
		return Item{}, err
	}
	return result, nil
}

func Update(db *mongo.Database, id string, newItem Item) (*mongo.UpdateResult, error) {
	collection := db.Collection("items")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": newItem}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	return updateResult, nil
}

func Delete(db *mongo.Database, id string) (*mongo.DeleteResult, error) {
	collection := db.Collection("items")
	filter := bson.M{"_id": id}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	return deleteResult, nil
}
