package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authorizer interface {
	Login(user, password string) (string, error)
	Validate(tokenString string) (*CustomClaims, error)
}

type Auth struct {
	adminUser     string
	adminPassword string
	tokenTTL      time.Duration
	secretKey     []byte
}

type CustomClaims struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func New(adminUser, adminPassword string, ttl time.Duration) *Auth {
	return &Auth{
		adminUser:     adminUser,
		adminPassword: adminPassword,
		tokenTTL:      ttl,
		secretKey: []byte("super-secret-key-for-task-7"),
	}
}

func (a *Auth) Login(user, password string) (string, error) {
	if user != a.adminUser || password != a.adminPassword {
		return "", fmt.Errorf("invalid credentials")
	}

	claims := CustomClaims{
		User: "superuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "superuser",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}

func (a *Auth) Validate(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
