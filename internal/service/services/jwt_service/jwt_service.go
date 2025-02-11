package jwt_service

import (
	"errors"
	"example/user/hello/types"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const kaccessTokenLife = 7 * 24 * time.Hour
const krefreshTokenLife = 7 * 24 * time.Hour

type JwtServiceImpl struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewJwtService() (*JwtServiceImpl, error) {
	accessSecretFromEnv := os.Getenv("ACCESS_SECRET")
	refreshSecretFromEnv := os.Getenv("REFRESH_SECRET")

	if accessSecretFromEnv == "" || refreshSecretFromEnv == "" {
		return nil, errors.New("missing JWT secrets in environment")
	}

	return &JwtServiceImpl{
		accessSecret:  []byte(accessSecretFromEnv),
		refreshSecret: []byte(refreshSecretFromEnv),
	}, nil
}

func (j *JwtServiceImpl) GenerateTokens(userID uint) (types.JwtPayload, error) {
	accessToken, err := j.createToken(userID, j.accessSecret, kaccessTokenLife)
	if err != nil {
		return types.JwtPayload{}, err
	}

	refreshToken, err := j.createToken(userID, j.refreshSecret, krefreshTokenLife)
	if err != nil {
		return types.JwtPayload{}, err
	}

	return types.JwtPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JwtServiceImpl) ValidateAccessToken(tokenString string) (types.Claims, error) {
	return j.parseToken(tokenString, j.accessSecret)
}

func (j *JwtServiceImpl) ValidateRefreshToken(tokenString string) (types.Claims, error) {
	return j.parseToken(tokenString, j.refreshSecret)
}

func (j *JwtServiceImpl) createToken(userID uint, secret []byte, duration time.Duration) (string, error) {
	claims := types.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (j *JwtServiceImpl) parseToken(tokenString string, secret []byte) (types.Claims, error) {
	claims := &types.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return types.Claims{}, err
	}

	if !token.Valid {
		return types.Claims{}, errors.New("invalid token")
	}

	return *claims, nil
}
