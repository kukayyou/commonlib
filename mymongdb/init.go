package mymongdb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Database struct {
	Mongo *mongo.Client
}

var DB *Database

//初始化
func Init(mongoUrl string) {
	DB = &Database{
		Mongo: SetConnect(mongoUrl),
	}
}

//连接设置
func SetConnect(mongoUrl string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 连接池
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl).SetMaxPoolSize(20))
	if err != nil {
		log.Println(err)
	}
	return client
}

type mgo struct {
	database   string
	collection string
}

func NewMgo(database, collection string) *mgo {
	return &mgo{
		database,
		collection,
	}
}
