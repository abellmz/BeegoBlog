package util

import (
	"context"
	"github.com/astaxie/beego/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var UserCollection *mongo.Collection

func InitMongo() {
	var err error
	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:28017"))
	if err != nil {
		logs.Error("Failed to connect to MongoDB: %v", err)
		return
	}

	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		logs.Error("Failed to connect to MongoDB: %v", err)
		return
	}
	UserCollection = Client.Database("beegodb").Collection("users")
	logs.Info("MongoDB connected successfully!")
}
