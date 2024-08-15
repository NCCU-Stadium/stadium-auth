package helper

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userID = uuid.New().String()

const (
	URI  = "/* mongodb+srv://... */"
	NAME = "CitizenMediaAuth"
	COLL = "RefreshTokens"
)

func TNewRefreshDB(uri, name, coll string) *RefreshHelper {
	client, err := mongo.Connect(context.TODO(), options.Client(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Panic(err)
		return nil
	}
	collection := client.Database(name).Collection(coll)
	err = newExpireIndex(collection) // Create index for expireAt field
	return &RefreshHelper{
		collection: collection,
	}
}

func TestGetRefreshMeta(t *testing.T) {
	rh := TNewRefreshDB(URI, NAME, COLL)
	if rh == nil {
		t.Error("Error while creating refresh helper")
		return
	}

	// 1. Save refresh meta
	tokenid, err := rh.SaveRefreshMeta(userID, time.Minute*1)
	if err != nil {
		t.Error(err)
		t.Error("Error while saving refresh meta")
	}

	// 2. Get refresh meta
	time.Sleep(time.Second * 5)
	_, err = rh.GetRefreshMeta(tokenid)
	if err != nil {
		t.Error(err)
		t.Error("Error while getting refresh meta")
	}

	// 3. Get refresh meta again (expect ErrorTokenUsed)
	time.Sleep(time.Second * 5)
	_, err = rh.GetRefreshMeta(tokenid)
	if err != ErrorTokenUsed {
		t.Error(err)
		t.Error("Error while getting refresh meta")
	}
}
