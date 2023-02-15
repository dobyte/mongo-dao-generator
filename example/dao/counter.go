package dao

import (
	"github.com/dobyte/mongo-dao-generator/example/dao/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type Counter struct {
	*internal.Counter
}

func NewCounter(db *mongo.Database) *Counter {
	return &Counter{Counter: internal.NewCounter(db)}
}
