package dao

import (
	"github.com/dobyte/mongo-dao-generator/example/dao/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserColumns = internal.UserColumns

type User struct {
	*internal.User
}

func NewUser(db *mongo.Database) *User {
	return &User{User: internal.NewUser(db)}
}
