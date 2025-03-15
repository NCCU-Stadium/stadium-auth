package helper

import (
	"auth-service/internal/config"
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshMeta struct {
	used     bool      `bson:"used"`
	UserMail string    `bson:"umail"`
	tokenID  string    `bson:"_id"`
	expireAt time.Time `bson:"expireAt"`
}

type RefreshHelper struct {
	collection *mongo.Collection
}

func NewRefreshDB(config *config.Config) *RefreshHelper {
	client, err := mongo.Connect(context.TODO(), options.Client(), options.Client().ApplyURI(config.RefreshDBURI))
	if err != nil {
		panic(err)
	}
	// check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Panic(err)
		return nil
	}
	collection := client.Database(config.RefreshDBName).Collection(config.RefreshDBCollection)
	err = newExpireIndex(collection) // Create index for expireAt field
	return &RefreshHelper{
		collection: collection,
	}
}

func (rh *RefreshHelper) Close() {
	rh.collection.Database().Client().Disconnect(context.Background())
	return
}

func newExpireIndex(collection *mongo.Collection) error {
	// Create index for expireAt field
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"expireAt": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	if err != nil {
		return err
	}
	return nil
}

// Define error types
var ErrorTokenNotFound error = errors.New("Token not found")
var ErrorTokenUsed error = errors.New("Token already used")

func (rh *RefreshHelper) GetRefreshMeta(tokenid string) (*RefreshMeta, error) {
	// Find the token by tokenid in database
	var result bson.M
	err := rh.collection.FindOne(context.TODO(), bson.M{"_id": tokenid}).Decode(&result)

	// If not found, return error
	if err == mongo.ErrNoDocuments {
		return nil, ErrorTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	// If found & not used, mark it as used and return refresh metadata
	if !result["used"].(bool) {
		_, err = rh.collection.UpdateOne(context.TODO(), bson.M{"_id": tokenid}, bson.M{"$set": bson.M{"used": true}})
		if err != nil {
			return nil, err
		}
		return &RefreshMeta{
			used:     false,
			UserMail: result["umail"].(string),
			tokenID:  result["_id"].(string),
		}, nil
	}

	return nil, ErrorTokenUsed
}

// Save the refresh metaData to database and return token id
func (rh *RefreshHelper) SaveRefreshMeta(userid string, expireAfter time.Duration) (string, error) {
	tokenid := uuid.New().String()
	_, err := rh.collection.InsertOne(context.TODO(), bson.M{
		"_id":      tokenid,
		"umail":    userid,
		"used":     false,
		"expireAt": time.Now().Add(expireAfter),
	})
	if err != nil {
		return "", err
	}
	return tokenid, nil
}
