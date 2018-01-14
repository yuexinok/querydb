package querydb

import (
	"math/rand"
	"strconv"
	"time"
)

func GetBranchName(name string, modnum int, mod string) string {
	if mod == "" {
		return name
	}
	num, err := strconv.Atoi(mod)
	if err != nil {
		dblog.Fatal("DB配置文件错误")
	}
	m := modnum % num
	name += strconv.Itoa(m)
	return name
}

func GetReadNumByRand(max int) int {
	if max < 1 {
		return 0
	}
	rand.Seed(time.Now().Unix())
	return rand.Intn(max)
}
