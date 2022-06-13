package tokenx

import (
	"fmt"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/warthecatalyst/douyin/rdb"
)

const (
	userId   = "user_id"
	username = "user_name"
	InvalidUserId = -1
)

func ParseToken(tokenJson string) (int64, string) {
	if tokenJson == "" {
		return InvalidUserId, ""
	}
	salts := rdb.GetAllSalts()
	result := make(chan jwt.MapClaims, 1)
	var wg sync.WaitGroup
	for index := 0; index < len(salts); index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			token, err := jwt.Parse(tokenJson, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(salts[index]), nil
			})
			if err == nil && token.Valid {
				claims, _ := token.Claims.(jwt.MapClaims)
				result <- claims

			}
		}(index)
	}
	wg.Wait()
	for item := range result {
		if item != nil {
			return int64(item[userId].(float64)), item[username].(string)
		}
	}

	return InvalidUserId, ""
}

func CreateToken(uid int64, uname string) string {
	mapClaims := jwt.MapClaims{
		userId:   uid,
		username: uname,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenJson, _ := token.SignedString(rdb.GetRandomSalt())
	return tokenJson
}

func CheckToken(userId int64, token string) bool {
	userIdFromToken, _ := ParseToken(token)
	if userId != userIdFromToken {
		return false
	}
	return true
}
