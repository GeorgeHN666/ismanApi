package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (app *application) Connect() *mongo.Client {

	op := options.Client().ApplyURI(app.config.db.uri)

	c, err := mongo.Connect(context.TODO(), op)
	if err != nil {
		app.errorLog.Panic(err)
	}

	err = c.Ping(context.TODO(), nil)
	if err != nil {
		app.errorLog.Panic(err)
	}

	return c
}

func (app *application) GetSubscriber(em string) (*Subscriber, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := app.Connect().Database(app.config.db.database).Collection("boletin")

	filter := bson.M{
		"email": bson.M{"$eq": em},
	}

	var res *Subscriber

	err := db.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (app *application) GetSubscriberList() ([]*Subscriber, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := app.Connect().Database(app.config.db.database).Collection("boletin")

	var res []*Subscriber

	filter := bson.M{}

	cursor, err := db.Find(ctx, filter)
	if err != nil {
		return res, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {

		var sub *Subscriber

		err := cursor.Decode(&sub)
		if err != nil {
			return res, err
		}

		res = append(res, sub)
	}

	if err := cursor.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (app *application) InsertSubscriber(s *Subscriber) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := app.Connect().Database(app.config.db.database).Collection("boletin")

	s.ID = primitive.NewObjectID()
	s.CreatedAt = time.Now()
	s.Status = 1
	code, _ := app.GenerateCode()
	s.Code = code

	_, err := db.InsertOne(ctx, s)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) DeleteSubscription(s *Subscriber) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := app.Connect().Database(app.config.db.database).Collection("boletin")

	filter := bson.M{
		"_id": bson.M{"$eq": s.ID},
	}

	_, err := db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
