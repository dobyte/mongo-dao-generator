package mail

import (
	"example/dao/mail/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type Columns = internal.Columns

type Mail struct {
	*internal.Mail
}

func NewMail(db *mongo.Database) *Mail {
	return &Mail{Mail: internal.NewMail(db)}
}
