package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID primitive.ObjectID `bson:"userid" json:"userid"`
	RoomID primitive.ObjectID `bson:"roomid" json:"roomid"`
	HotelID primitive.ObjectID `bson:"hotelid" json:"hotelid"`
	Start time.Time `bson:"start" json:"start"`
	End time.Time `bson:"end" json:"end"`
}

type BookingParams struct {
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
	NumPerson int `json:"num"`
	Validator
}

func (b *BookingParams) Validate(rt RoomType) bool{
	b.CheckField(CheckDates(b.Start, b.End), "date", "Dates must be correctly chosen")
	if b.NumPerson == 0 {
		b.AddFieldError("num", "Number of persons must be bigger than 0")
	}

	if rt == Single && b.NumPerson > 1{
		b.AddFieldError("num", "Only one person can be allowed to the single room")
	}else if rt == Double && b.NumPerson > 2{
		b.AddFieldError("num", "Up to two person can be allowed to the double room")
	}else if rt == Deluxe && b.NumPerson > 4{
		b.AddFieldError("num", "Up to four person can be allowed to the deluxe room")
	}

	return b.Valid()
}
