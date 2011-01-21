package morestore

import (
	"./statmsg"
	"github.com/mikejs/gomongo/mongo"
	"redis"
	"os"
	"fmt"
)

type Context struct {
	redisClient redis.Client

	mongoPool chan *mongo.Collection
	mongoAddr string
	mongoDb   string
}

func (c *Context) popMongoCollection() (*mongo.Collection, os.Error) {
	collection := <- c.mongoPool

	if collection == nil {
		conn, err := mongo.Connect(c.mongoAddr)
		if err != nil {
			return nil, err
		}
		collection = conn.GetDB(c.mongoDb).GetCollection("logs")
	}

	return collection, nil
}

func (c *Context) pushMongoCollection(collection *mongo.Collection) {
	c.mongoPool <- collection
}

func (c *Context) Update(stat *statmsg.Statmsg) {
	_, err := c.redisClient.Incr("hit:" + stat.Key)
	if err != nil {
		fmt.Printf("Error from redis: %s\n", err.String())
	}

	collection, err := c.popMongoCollection()
	if err != nil {
		fmt.Printf("Error from mongo: %s\n", err.String())
	}
	doc, _ := mongo.Marshal(stat)
	collection.Insert(doc)
	c.pushMongoCollection(collection)
}

func Setup(mongoAddr string, mongoDb string,
	redisAddr string, redisDb int, poolSize int) (c *Context) {
	c = new(Context)
	c.redisClient.Addr = redisAddr
	c.redisClient.Db = redisDb
	c.redisClient.MaxPoolSize = poolSize

	c.mongoPool = make(chan *mongo.Collection, poolSize)
	for i := 0; i < poolSize; i++ {
		c.mongoPool <- nil
	}
	c.mongoAddr = mongoAddr
	c.mongoDb = mongoDb

	return c
}
