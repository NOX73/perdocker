package perd

import "math/rand"
import "strconv"

const (
  randomIdLength = 10
)

func randomId () string {
  b := make([]byte, randomIdLength)

  for i := range b {
    b[i] = byte(rand.Int63() & 0xff)
  }

  return string(b)
}

var fileNameId int64
func uniqFileName () string {
  fileNameId++
  return strconv.FormatInt(fileNameId, 10)
}
