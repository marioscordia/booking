package api

import (
	"booking/db"
	"booking/types"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	DB *db.DB
}

func NewRoomHandler(db *db.DB) *RoomHandler {
	return &RoomHandler{
		DB: db,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	info := c.Context().UserValue("info").(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(info["id"].(string))
	if err != nil {
		fmt.Println(err)
		return err
	}

	var params types.BookingParams
	if err := c.BodyParser(&params); err != nil {
		fmt.Println(err)
		return err
	}
	
	roomID, err := primitive.ObjectIDFromHex(c.Params("roomid"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	room, err := h.DB.GetRoomByID(c.Context(), c.Params("roomid"))
	if err != nil {
		if errors.Is(err, db.ErrNoRecord){
			return NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		}
		fmt.Println(err)
		return err
	}

	if !params.Validate(room.Type){
		return c.Status(http.StatusBadRequest).JSON(params.FieldErrors)
	}

	filter := bson.M{
		"roomid": roomID,
		"start": bson.M{
			"$gte": params.Start,
		},
		"end": bson.M{
			"$lte": params.End,
		},
	}
	
	bookings, err := h.DB.GetBookings(c.Context(), filter)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if len(bookings) > 0 {
		return c.Status(http.StatusBadRequest).SendString("The room is already booked in this date range")
	}

	book := &types.Booking{
		UserID: userID,
		RoomID: room.ID,
		HotelID: room.HotelID,
		Start: params.Start,
		End: params.End,
	}
	
	err = h.DB.BookRoom(c.Context(), book)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.SendString("You booked successfully!")
}

func (h *RoomHandler) HandleGetBookings(c *fiber.Ctx) error {
	info := c.Context().UserValue("info").(jwt.MapClaims)
	
	userID, err := primitive.ObjectIDFromHex(info["id"].(string))
	if err != nil {
		fmt.Println(err)
		return err
	}

	filter := bson.M{
		"userid": userID,
	}

	bookings, err := h.DB.GetBookings(c.Context(), filter)
	if err != nil {
		if errors.Is(err, db.ErrNoRecord){
			return c.SendString("You have no bookings yet!")
		}
		fmt.Println(err)
		return err
	}
	if len(bookings) == 0 {
		return c.SendString("You have no bookings yet!")
	}

	return c.JSON(bookings)	
}

func (h *RoomHandler) HandleCancelBooking(c *fiber.Ctx) error{
	info := c.Context().UserValue("info").(jwt.MapClaims)

	userID, err := primitive.ObjectIDFromHex(info["id"].(string))
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	id := c.Params("id")
	bID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound)) 
	}

	err = h.DB.CancelBooking(c.Context(), bson.M{"_id": bID, "userid": userID})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.SendString("Canceled booking successfully!")
}
