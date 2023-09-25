package api

import (
	"booking/db"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	DB *db.DB
}

func NewHotelHandler(db *db.DB) *HotelHandler {
	return &HotelHandler{
		DB: db,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error{
	hotels, err := h.DB.GetHotels(c.Context())
	if err != nil{
		if errors.Is(err, db.ErrNoRecord){
			return NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		}
		ErrorLog(err)
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error{
	id := c.Params("id")
	hotel, err := h.DB.GetHotelByID(c.Context(), id)
	if err != nil{
		if errors.Is(err, db.ErrNoRecord){
			return NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		}
		ErrorLog(err)
		return err
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRoomsByID(c *fiber.Ctx) error{
	id := c.Params("id")
	rooms, err := h.DB.GetRoomsByHotelID(c.Context(), id)
	if err != nil{
		if errors.Is(err, db.ErrNoRecord){
			return NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		}
		ErrorLog(err)
		return err
	}

	return c.JSON(rooms)
}