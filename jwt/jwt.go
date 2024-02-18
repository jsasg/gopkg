package jwt

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenExpired = jwt.ErrTokenExpired

type Claims struct {
	UID int64 `json:"uid"`
	jwt.RegisteredClaims
}

type JWT struct {
	Key string
}

func New(key string) *JWT {
	return &JWT{
		Key: key,
	}
}

// AccessToken 生成授权用的 access_token
func (j *JWT) AccessToken(jwtID string, uid int64, expire int64) (string, int64, error) {
	var (
		token     string
		issuedAt  = time.Now()
		expiresAt = issuedAt.Add(time.Duration(expire) * time.Second)
		err       error
	)

	claims := Claims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        jwtID,
		},
	}

	if token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.Key)); err != nil {
		return "", 0, err
	}

	return token, expiresAt.Unix(), nil
}

// RefreshToken 生成刷新授权用的 refresh_token
func (j *JWT) RefreshToken(jwtID string, expire int64) (string, int64, error) {
	var (
		token     string
		issuedAt  = time.Now()
		expiresAt = issuedAt.Add(time.Duration(expire) * time.Second)
		err       error
	)

	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        jwtID,
	}

	if token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.Key)); err != nil {
		return "", 0, err
	}

	return token, expiresAt.Unix(), nil
}

// Parse 解析 token
func (j *JWT) Parse(tokenString string) (*Claims, error) {
	var (
		claims Claims
		err    error
	)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Key), nil
	})

	return &claims, err
}
