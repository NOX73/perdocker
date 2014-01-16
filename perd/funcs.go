package perd

import (
  "math/rand"
  "strconv"
  "sync"
)

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
var fileNameIdMutex sync.Mutex

func uniqFileName () string {
  fileNameIdMutex.Lock()
    fileNameId++
    id := fileNameId
  fileNameIdMutex.Unlock()

  return strconv.FormatInt(id, 10)
}
