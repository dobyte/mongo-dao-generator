package main

import (
	"context"
	"github.com/dobyte/mongo-dao-generator/example/dao"
	"github.com/dobyte/mongo-dao-generator/example/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func main() {
	var (
        uri     = "mongodb://root:12345678@127.0.0.1:27017"
		opts    = options.Client().ApplyURI(uri)
		baseCtx = context.Background()
	)

    ctx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
    client, err := mongo.Connect(ctx, opts)
    cancel()
    if err != nil {
        log.Fatalf("connect mongo server failed: %v", err)
    }

	ctx, cancel = context.WithTimeout(baseCtx, 5*time.Second)
    defer cancel()
	err = client.Ping(ctx, readpref.Primary())
    cancel()
	if err != nil {
		log.Fatalf("ping mongo server failed: %v", err)
	}

	db := client.Database("dao_test")

	mailDao := dao.NewMail(db)

	_, err = mailDao.InsertOne(baseCtx, &model.Mail{
		Title:    "mongo-dao-generator introduction",
		Content:  "the mongo-dao-generator is a tool for automatically generating MongoDB Data Access Object.",
		Sender:   1,
		Receiver: 2,
		Status:   1,
	})
	if err != nil {
		log.Fatalf("failed to insert into mongo database: %v", err)
	}

	mail, err := mailDao.FindOne(baseCtx, func(cols *dao.MailColumns) interface{} {
		return bson.M{cols.Receiver: 2}
	})
	if err != nil {
		log.Fatalf("failed to find a row of data from mongo database: %v", err)
	}

	log.Printf("%+v", mail)
}
