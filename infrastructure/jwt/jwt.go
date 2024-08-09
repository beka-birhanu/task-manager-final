package jwt

import (
	"errors"
	"time"

	usermodel "github.com/beka-birhanu/models/user"
	"github.com/dgrijalva/jwt-go"
)

// Service implements the ijwt.IService interface for handling JWT operations.
type Service struct {
	secretKey string
	issuer    string
	expTime   time.Duration
}

// Config holds the configuration for creating a new JWT Service.
type Config struct {
	SecretKey string
	Issuer    string
	ExpTime   time.Duration
}

func New(config Config) *Service {
	return &Service{
		secretKey: config.SecretKey,
		issuer:    config.Issuer,
		expTime:   config.ExpTime,
	}
}

func (s *Service) Generate(user *usermodel.User) (string, error) {
	expirationTime := time.Now().UTC().Add(s.expTime).Unix()
	claims := jwt.MapClaims{
		"user_id":  user.ID().String(),
		"is_admin": user.IsAdmin(),
		"exp":      expirationTime,
		"iss":      s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *Service) Decode(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, s.getSigningKey)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *Service) getSigningKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return []byte(s.secretKey), nil
}
