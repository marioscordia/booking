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


type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, bson.M, bson.M) error
	GetHotels(context.Context) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	// dbname string
	coll *mongo.Collection 
}

func NewMongoHotelStore(client *mongo.Client, config *config.Config) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll: client.Database(config.DB.Name).Collection(config.DB.Hotel),
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	
	hotel.ID = res.InsertedID.(primitive.ObjectID)

	return hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	var hotels []*types.Hotel

	curr , err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer curr.Close(ctx)

	if err = curr.All(ctx, &hotels); err != nil {
		if errors.Is(err, ErrNoRecord){
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return hotels, nil
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	hotelID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel *types.Hotel
	err = s.coll.FindOne(ctx, bson.M{"_id": hotelID}).Decode(&hotel)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return hotel, nil
}