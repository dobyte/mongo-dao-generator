package dao

import (
	"example/dao/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mail struct {
	*internal.Mail
}

func NewMail(db *mongo.Database) *Mail {
	return &Mail{Mail: internal.NewMail(db)}
}
