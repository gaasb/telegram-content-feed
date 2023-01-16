package bot

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE_NAME = "DATABASE_NAME"

	FEED_COLLECTION             = "feed"
	POSTS_COLLECTION            = "posts"
	USERS_WITH_ROLE_COLLECTION  = "users"
	TAGS_COLLECTION             = "tags"
	MEDIA_FOR_REVIEW_COLLECTION = "media_for_review"
)

var (
	database *mongo.Database
)

func GetById() {
	var result interface{}
	entities, _ := database.Collection(FEED_COLLECTION).Find(context.TODO(), bson.D{})
	entities.All(context.TODO(), &result)
	fmt.Println(result)
}
func AddMedia(msg interface{}) error {
	//database.Aggregate() TODO CHECK IS IN FEED AND media_for_review COLLECTION!
	_, err := database.Collection(MEDIA_FOR_REVIEW_COLLECTION).InsertOne(context.TODO(), msg)
	return err
}
func AddMediaToFeed(msg interface{}) error {
	_, err := database.Collection(FEED_COLLECTION).InsertOne(context.TODO(), msg)
	return err
}
func FindAllMedia() []*MediaMessage {
	var res []*MediaMessage
	result, err := database.Collection(MEDIA_FOR_REVIEW_COLLECTION).Find(context.TODO(), bson.D{})
	if err != nil {
		return nil
	}
	result.All(context.TODO(), &res)
	return res
}
func RemoveMediaByID(id interface{}) error {
	_, err := database.Collection(MEDIA_FOR_REVIEW_COLLECTION).DeleteOne(context.TODO(), bson.D{{"_id", id}})
	return err
}

func instanceDatabaseCollections() {
	_ = database.CreateCollection(context.TODO(), FEED_COLLECTION)
	_ = database.CreateCollection(context.TODO(), POSTS_COLLECTION)
	_ = database.CreateCollection(context.TODO(), USERS_WITH_ROLE_COLLECTION)
	_ = database.CreateCollection(context.TODO(), MEDIA_FOR_REVIEW_COLLECTION)
}
