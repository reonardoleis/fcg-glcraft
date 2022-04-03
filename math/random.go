package math2

import "math/rand"

func SetSeed(seed int64) {
	rand.Seed(seed)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
