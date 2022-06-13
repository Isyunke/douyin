package rdb

import "strconv"

func AddToken(userId int64, token string) error {
	return rdb.Set(strconv.FormatInt(userId, 10), token, 0).Err()
}

func GetToken(userId int64) string {
	return rdb.Get(strconv.FormatInt(userId, 10)).Val()
}
