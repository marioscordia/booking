package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Location string `bson:"location" json:"location"`
	Rooms []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

type RoomType int

const (
	_ RoomType = iota
	Single 
	Double
	Deluxe
)

type BedSize int

const (
	_ BedSize = iota
	Small
	Medium
	Large
)

type Room struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type RoomType `bson:"type" json:"type"`
	Bed BedSize `bson:"bed" json:"bed"`
	SeaSide bool `bson:"sea_side" json:"sea_side"`
	Price float64 `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelid" json:"hotelid"`
}

func (r Room) CheckNum(n int) bool {
	if r.Type == Single{
		return n == 1
	}else if r.Type == Double{
		return n <= 2
	}

	return n <= 4
}