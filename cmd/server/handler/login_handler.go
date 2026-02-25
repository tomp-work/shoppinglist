package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

var (
	accessSecret  = []byte("ACCESS_SECRET")
	refreshSecret = []byte("REFRESH_SECRET")
)

type JwtCustomClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"userId"`
}

// generateToken generates a JWT.
func generateToken(userID string, secret []byte, ttl time.Duration) (string, error) {
	claims := JwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token in generateToken: %w", err)
	}
	return signed, nil
}

// login.
func (h *Handler) Login(c *echo.Context) error {
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	// User database!
	type User struct {
		Username string
		Password string
		ID       string
	}
	userDB := []User{
		{Username: "admin", Password: "admin1", ID: "111"},
		{Username: "tom", Password: "tom2", ID: "222"},
		{Username: "clare", Password: "clare3", ID: "333"},
	}

	user := (*User)(nil)
	for _, u := range userDB {
		if req.Username == u.Username && req.Password == u.Password {
			user = &u
			break
		}
	}
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	accessToken, err := generateToken(user.ID, accessSecret, 15*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to generate access token in handler.Login(): %w", err)
	}
	refreshToken, err := generateToken(user.ID, refreshSecret, 7*24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token in handler.Login(): %w", err)
	}

	rsp := struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return c.JSON(http.StatusOK, &rsp)
}
