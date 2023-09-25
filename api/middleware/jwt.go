package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTauthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthenticated") 
	}

	claims, err := ParseToken(token) 
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("unauthenticated")
	}
	//  fmt.Println(claims)
	exp := claims["exp"].(float64)
	if float64(time.Now().Unix()) >= exp {
		fmt.Println("token is expired")
		return fmt.Errorf("unauthorized")
	}
	fmt.Println("authorized")
	c.Context().SetUserValue("info", claims)
	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	info := c.Context().UserValue("info").(jwt.MapClaims)
	IsAdmin := info["admin"].(bool)
	if !IsAdmin{
		fmt.Println("Only admin is allowed")
		return fmt.Errorf("do not have access")
	}
	c.Context().SetUserValue("info", info)
	return c.Next()
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid{
		return nil, fmt.Errorf("token is invalid")
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok{
		// fmt.Println(claims)
		return claims, nil
	}

	return nil, fmt.Errorf("unauthenticated")
	
}