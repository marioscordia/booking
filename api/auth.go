package api

import (
	"booking/db"
	"booking/types"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	DB *db.DB
}

func NewAuthHandler(db *db.DB) *AuthHandler {
	return &AuthHandler{
		DB: db,
	}
}

type AuthResponse struct {
	User *types.User `json:"user"`
	Token string `json:"token"`
}
 
func (h *AuthHandler) Login(c *fiber.Ctx)	error{
	var LoginParams types.LoginParams

	if err := c.BodyParser(&LoginParams); err != nil{
		fmt.Println(err)
		return err
	}

	if !LoginParams.Validate(){
		return c.JSON(LoginParams.FieldErrors)
	}

	user, err := h.DB.Authenticate(c.Context(), &LoginParams)
	if err != nil{
		if errors.Is(err, db.ErrInvalidCredentials){
			return NewError(fiber.StatusBadRequest, "Invalid credentials")
		}
		fmt.Println(err)
		return err
	}

	token, err := CreateToken(user)
	if err != nil{
		fmt.Println(err)
		return err
	}
	
	resp := AuthResponse{
		User: user,
		Token: token,
	}

	return c.JSON(resp)
}

func CreateToken(u *types.User) (string, error) {
	claims := jwt.MapClaims{
		"id":u.ID,
		"email":u.Email,
		"admin": u.IsAdmin,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(key))
	if err != nil{
		return "", err
	}

	return tokenStr, nil
}