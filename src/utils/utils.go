package utils

import (
	"DHT/src/config"
	"crypto/sha1"
	"encoding/binary"
	"math"
)

func Hash(input string) int {
	h := sha1.New()
	h.Write([]byte(input))
	bs := h.Sum(nil)
	return int((binary.BigEndian.Uint64(bs)) % uint64(math.Pow(2, config.M)))
}

func BuildPath(id int) []int {
	path := []int{}
	for id > 0 {
		path = append([]int{id}, path...)
		id = (id - 1) / 2
	}
	return append([]int{0}, path...)
}
