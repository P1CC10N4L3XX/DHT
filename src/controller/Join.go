package controller

import (
	"DHT/src/config"
	"DHT/src/utils"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"time"
)

type Join struct{}

func (j *Join) InitConnection(port string) int {
	hostname, _ := os.Hostname()
	ts := time.Now().UnixNano()
	meta := fmt.Sprintf("%s-%s-%d", hostname, port, ts)
	id := utils.Hash(meta)
	path := utils.BuildPath(id)

}
