package database

import (
	"context"
	"log"

	"github.com/Youn-sik/KSCONNECT_ADMIN/setting"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongodbConnection() (*mongo.Client, error) {
	config, err := setting.LoadConfigSettingJSON()

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	credential := options.Credential{
		Username: config.Mongodb.User,
		Password: config.Mongodb.Password,
	}
	clientOptions := options.Client().ApplyURI("mongodb://" + config.Mongodb.Host + ":" + config.Mongodb.Port).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("[SERVER] => MongoDB connection made")
	return client, err
}
