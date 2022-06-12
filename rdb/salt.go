package rdb

const (
	keySalt = "salt"
)

func GetAllSalts() []string {
	return rdb.SMembers(keySalt).Val()
}

func GetRandomSalt() []byte {
	return []byte(rdb.SRandMemberN(keySalt, 1).Val()[0])
}
