package user

import (
	"example/dao/user/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type Columns = internal.Columns

type User struct {
	*internal.User
}

func NewUser(db *mongo.Database) *User {
	return &User{User: internal.NewUser(db)}
}
