package db

import (
	"context"
	"crawler/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListLinks() (list []model.VisitedLink, err error) {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection("links")

	opts := options.Find().SetSort(bson.D{{"visited_date", -1}})
	cursor, err := c.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		return
	}

	err = cursor.All(context.TODO(), &list)

	return

	return
}
