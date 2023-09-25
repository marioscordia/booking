package db

import (
	"booking/config"
	"booking/types"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookStore interface{
	BookRoom(context.Context, *types.Booking) error
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	CancelBooking(context.Context, bson.M) error
	DeleteBookings(context.Context, string) error
}

type MongoBookStore struct {
	client *mongo.Client
	// dbname string
	coll *mongo.Collection 
}

func NewMongoBookStore(client *mongo.Client, config *config.Config) *MongoBookStore {
	return &MongoBookStore{
		client: client,
		coll: client.Database(config.DB.Name).Collection(config.DB.Book),
	}
}

func (s *MongoBookStore) BookRoom(ctx context.Context, book *types.Booking) error {
	_, err := s.coll.InsertOne(ctx, book)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *MongoBookStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	
	var bookings []*types.Booking
	res, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer res.Close(ctx)
	
	if err := res.All(ctx, &bookings); err != nil {
		return nil, err
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookStore) CancelBooking(ctx context.Context, filter bson.M) error {

	_, err := s.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *MongoBookStore) DeleteBookings(ctx context.Context, userid string) error {
	userID, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteMany(ctx, bson.M{"userid": userID})
	if err != nil{
		if !errors.Is(err, mongo.ErrNoDocuments){
			return err
		}
	}	

	return nil
}