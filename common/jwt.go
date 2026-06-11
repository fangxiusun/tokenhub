package common

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

// JWTSecret is the secret key for JWT signing
var JWTSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Generate random secret if not set
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			log.Fatal("Failed to generate JWT secret:", err)
		}
		JWTSecret = []byte(hex.EncodeToString(b))
		log.Println("WARNING: JWT_SECRET not set, generated random secret. Set JWT_SECRET environment variable for production.")
	} else {
		JWTSecret = []byte(secret)
	}
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserId    int    `json:"user_id"`
	Username  string `json:"username"`
	Role      int    `json:"role"`
	Status    int    `json:"status"`
	GroupId   string `json:"group_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(userId int, username string, role, status int, groupId string) (string, error) {
	claims := JWTClaims{
		UserId:   userId,
		Username: username,
		Role:     role,
		Status:   status,
		GroupId:  groupId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "tokenhub",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ParseJWT parses and validates a JWT token
func ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return JWTSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	return ParseJWT(tokenString)
}

