package db

import (
	"booking/config"
	"booking/types"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) error
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, *types.UpdatedUser) error
	Authenticate(context.Context, *types.LoginParams) (*types.User, error)
}

type MongoUserStore struct{
	client *mongo.Client
	// dbname string
	coll *mongo.Collection 
}

func NewMongoUserStore(client *mongo.Client, config *config.Config) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll: client.Database(config.DB.Name).Collection(config.DB.User),
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNoRecord
	}

	var user *types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(user); err != nil{
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	var users []*types.User

	curr, err := s.coll.Find(ctx, bson.M{"isAdmin": false})
	if err != nil {
		return nil, err
	}

	if err = curr.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, u *types.User) error {
	_, err := s.coll.InsertOne(ctx, u)
	if err != nil {
		if mongo.IsDuplicateKeyError(err){
			return ErrDuplicateEmail
		} else {
			return err
		}
	}

	return  nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrNoRecord
	}
	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNoRecord
	}
	
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, params *types.UpdatedUser) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": params})
	if err != nil {
		if mongo.IsDuplicateKeyError(err){
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (s *MongoUserStore) Authenticate(ctx context.Context, params *types.LoginParams) (*types.User, error) {
	
	filter := bson.M{
		"email": params.Email,
	}

	var user types.User

	err := s.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	
	return &user, nil
}




 