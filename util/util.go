package util

import (
	"github.com/google/uuid"
	"github.com/ser163/WordBot/generate"
	"strings"
)

func CreateUuid() int64 {
	return int64(uuid.New().ID())
}

func CreateRandomString(count int) []string {
	var randStrs []string
	for i := 1; i <= count; i++ {
		wordList, _ := generate.GenRandomMix(10)
		randStrs = append(randStrs, wordList.Word)
	}
	return randStrs
}

func StringStrip(input string) string {
	if input == "" {
		return ""
	}
	return strings.Join(strings.Fields(input), "")
}
