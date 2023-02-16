package dao

import (
	"github.com/dobyte/mongo-dao-generator/example/dao/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type MailColumns = internal.MailColumns

type Mail struct {
	*internal.Mail
}

func NewMail(db *mongo.Database) *Mail {
	return &Mail{Mail: internal.NewMail(db)}
}
