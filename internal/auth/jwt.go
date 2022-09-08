package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"tinderutf/utils/env"
)

type jwtService struct {
	secretKey string
	issure    string
}

var JWT = NewJWTService()

func NewJWTService() *jwtService {
	return &jwtService{
		secretKey: env.GetEnv("JWT_SECRET_KEY", "ashudaudhuas1231sadas"),
		issure:    "tinder-issuer",
	}
}

func (s *jwtService) GenerateToken(userId string) (string, error) {
	claim := &jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iss": s.issure,
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *jwtService) ValidateToken(token string) bool {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token: %v", token)
		}

		return []byte(s.secretKey), nil
	})

	return err == nil
}

func (s *jwtService) GetUserDataFromToken(t string) (UserData, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid Token: %v", t)
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return UserData{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["sub"].(string)
		if !ok {
			return UserData{}, fmt.Errorf("cannot parse")
		}

		return UserData{
			Id: id,
		}, nil
	}

	return UserData{}, err
}

type UserData struct {
	Id string
}
