package helper

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetAuthId(c echo.Context) (int, error) {
	claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)
	authID, err := strconv.Atoi(claims.ID)

	if err != nil {
		return -1, err
	}

	return authID, nil
}

func Login(userId int, name string) (map[string]interface{}, error) {
	// Default durasi token, misalnya 24 jam
	tokenDuration := 24 * time.Hour

	expiredTime := time.Now().Local().Add(tokenDuration)

	claims := JwtCustomClaims{
		ID: strconv.Itoa(userId),
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "skripsi",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token, errToken := NewTokenUseCase().GenerateAccessToken(claims)
	if errToken != nil {
		return nil, errors.New("internal server error")
	}

	return map[string]interface{}{
		"token":      token,
		"expires_at": expiredTime.Format(time.RFC3339),
	}, nil
}
