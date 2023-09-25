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


type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRoomsByHotelID(context.Context, string) ([]*types.Room, error)
	GetRoomByID(context.Context, string) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll *mongo.Collection 
}

func NewMongoRoomStore(client *mongo.Client, config *config.Config) *MongoRoomStore {
	return &MongoRoomStore{
		client: client,
		coll: client.Database(config.DB.Name).Collection(config.DB.Room),
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)
	return room, nil
}

func (s *MongoRoomStore) GetRoomsByHotelID(ctx context.Context, id string) ([]*types.Room, error) {
	hotelID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"hotelid": hotelID}
	var rooms []*types.Room
	res, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := res.All(ctx, &rooms); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) GetRoomByID(ctx context.Context, id string) (*types.Room, error) {
	roomID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": roomID}
	var room *types.Room
	err = s.coll.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return room, nil
}