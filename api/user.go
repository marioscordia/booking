package api

import (
	"booking/db"
	"booking/types"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	DB *db.DB
}

func NewUserHandler(db *db.DB) *UserHandler {
	return &UserHandler{
		DB: db,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	info := c.Context().UserValue("info").(jwt.MapClaims)
	userID := info["id"].(string)

	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		fmt.Println(err)
		return err
	}
	
	if !params.Validate(){
		return c.Status(fiber.StatusBadRequest).JSON(params.FieldErrors)
	}

	updated, err := types.NewUpdatedUser(params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	err = h.DB.UpdateUser(c.Context(), userID, updated)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.SendString("Info has been updated!")
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		fmt.Println(err)
		return err
	}

	if !params.Validate(){
		return c.Status(fiber.StatusBadRequest).JSON(params.FieldErrors)
	}

	user, err := types.NewUserParams(params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	err = h.DB.InsertUser(c.Context(), user)
	if err != nil {
		if errors.Is(err, db.ErrDuplicateEmail){
			return NewError(fiber.StatusBadRequest, "User with such email already exists")
		}
		fmt.Println(err)
		return err
	}

	return c.SendString("You have successfully registered!")
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error{
	id := c.Params("id")
	user, err := h.DB.GetUserById(c.Context(), id)
	if err != nil{
		if errors.Is(err, db.ErrNoRecord){
			return NewError(fiber.StatusNotFound, "User with such id does not exist")
		}
		fmt.Println(err)
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error{
	users, err := h.DB.GetUsers(c.Context())
	if err != nil{
		fmt.Println(err)
		return err
	}
	if len(users) == 0{
		c.SendString("There are no users yet!")
	}

	return c.JSON(users)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	err := h.DB.DeleteUser(c.Context(), userID)
	if err != nil {
		if errors.Is(err, db.ErrNoRecord){
			return NewError(fiber.StatusNotFound, "No user with such id")
		}
		fmt.Println(err)
		return err
	}

	err = h.DB.DeleteBookings(c.Context(), userID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.SendString("User successfully deleted!")
}