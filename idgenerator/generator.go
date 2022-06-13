package idgenerator

import "github.com/google/uuid"

func GenerateVid() int64 {
	return int64(uuid.New().ID())
}

func GenerateUid()int64{
	return int64(uuid.New().ID())
}