package user

import (
	"context"
	"example/dao/user/internal"
	"example/model/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	*internal.User
}

func NewUser(db *mongo.Database) *User {
	return &User{User: internal.NewUser(db)}
}

func (dao *User) FindOneByWechat(ctx context.Context, openID string) (*user.User, error) {
	return dao.FindOne(ctx, func(cols *internal.Columns) interface{} {
		return bson.M{cols.ThirdPlatforms: bson.M{
			"wechat": openID,
		}}
	})
}
