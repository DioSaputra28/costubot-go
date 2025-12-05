package utils

import (
	"contact-management/src/apps"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	rdb := apps.RedisClient()
	ctx := context.Background()
	userKey := fmt.Sprintf("user_token:%s", username)
	prevToken, err := rdb.Get(ctx, userKey).Result()
	if err == nil {
		prevTokenKey := fmt.Sprintf("token:%s", prevToken)
		rdb.Del(ctx, prevTokenKey)
	} else if err != redis.Nil {
		apps.LoggingApp().Warn("Failed to read previous token", err)
	}

	tokenKey := fmt.Sprintf("token:%s", token)
	if err := rdb.Set(ctx, tokenKey, username, 24*time.Hour).Err(); err != nil {
		apps.LoggingApp().Error("Failed to store token in Redis", err)
		return "", err
	}

	if err := rdb.Set(ctx, userKey, token, 24*time.Hour).Err(); err != nil {
		apps.LoggingApp().Error("Failed to update user token mapping", err)
		return "", err
	}

	return token, nil
}

func VerifyToken(tokenString string) (string, error) {
	rdb := apps.RedisClient()
	ctx := context.Background()

	tokenKey := fmt.Sprintf("token:%s", tokenString)
	username, err := rdb.Get(ctx, tokenKey).Result()
	if err != nil {
		return "", errors.New("Invalid token")
	}
	return username, nil
}

func RevokeToken(tokenString, username string) error {
	rdb := apps.RedisClient()
	ctx := context.Background()

	userKey := fmt.Sprintf("user_token:%s", username)
	tokenKey := fmt.Sprintf("token:%s", tokenString)

	if err := rdb.Del(ctx, userKey).Err(); err != nil {
		apps.LoggingApp().Error("Failed to revoke user token", err)
		return err
	}
	if err := rdb.Del(ctx, tokenKey).Err(); err != nil {
		apps.LoggingApp().Error("Failed to revoke token", err)
		return err
	}
	return nil
}
