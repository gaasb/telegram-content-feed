package bot

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE_NAME = "DATABASE_NAME"

	FEED_COLLECTION             = "feed"
	POSTS_COLLECTION            = "posts"
	TAGS_COLLECTION             = "tags"
	MEDIA_FOR_REVIEW_COLLECTION = "media_for_review"
	USERS_WITH_ROLES_COLLECTION = "administrators"
)

var (
	database *mongo.Database
)

func FindMediaById(id interface{}) *MediaMessage {
	var result *MediaMessage
	err := database.Collection(FEED_COLLECTION).FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return nil
	}
	return result
}
func AddMedia(msg interface{}) error {
	//database.Aggregate() TODO CHECK IS IN FEED AND media_for_review COLLECTION!
	_, err := database.Collection(MEDIA_FOR_REVIEW_COLLECTION).InsertOne(context.TODO(), msg)
	return err
}
func AddMediaToFeed(msg interface{}) error {
	_, err := database.Collection(FEED_COLLECTION).InsertOne(context.TODO(), msg)
	fmt.Println(err)
	return err
}
func FindFirstMedia() *MediaMessage {
	var res *MediaMessage
	_ = database.Collection(MEDIA_FOR_REVIEW_COLLECTION).FindOne(context.TODO(), bson.D{}).Decode(&res)
	return res

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
func InsertTag(tag *TagsStorage) error {
	_, err := database.Collection(TAGS_COLLECTION).InsertOne(context.TODO(), tag)
	return err
}
func GetTagsByTagType(value string) ([]*TagsStorage, error) {
	var result []*TagsStorage
	rawResult, err := database.Collection(TAGS_COLLECTION).Find(context.TODO(), bson.D{{"type", value}})
	if err != nil {
		return nil, err
	}
	rawResult.All(context.TODO(), &result)
	fmt.Println(result)
	return result, nil
}
func RemoveTagById(hexValue string) error {
	object, objectErr := primitive.ObjectIDFromHex(hexValue)
	if objectErr != nil {
		return errors.New("Invalid hex value for instance objectId in remove tag ")
	}
	_, err := database.Collection(TAGS_COLLECTION).DeleteOne(context.TODO(), bson.D{{"_id", object}})
	fmt.Println(err)
	return err
}
func InsertAdministrator(user *Administrator) error {
	if _, err := database.Collection(USERS_WITH_ROLES_COLLECTION).InsertOne(context.TODO(), user); err != nil {
		return err
	} else {
		return nil
	}
}
func FindAllAdministrators() map[int64][]UserType {
	var result = map[int64][]UserType{}
	var re []Administrator
	rawResult, err := database.Collection(USERS_WITH_ROLES_COLLECTION).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil
	}
	rawResult.All(context.TODO(), &re)
	for _, item := range re {
		result[item.UserID] = item.Rights
	}
	fmt.Println(result)
	//maps.Copy(result, re[0])
	//fmt.Println(result)
	return result
}

func instanceDatabaseCollections() {
	_ = database.CreateCollection(context.TODO(), FEED_COLLECTION)
	_ = database.CreateCollection(context.TODO(), POSTS_COLLECTION)
	_ = database.CreateCollection(context.TODO(), MEDIA_FOR_REVIEW_COLLECTION)
	_ = database.CreateCollection(context.TODO(), TAGS_COLLECTION)
	_ = database.CreateCollection(context.TODO(), USERS_WITH_ROLES_COLLECTION)
	createUniqueIndexForCaption()
}
