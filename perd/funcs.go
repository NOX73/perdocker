package perd

import (
	"strconv"
	"sync"
)

var fileNameID int64
var fileNameIDMutex sync.Mutex

func uniqFileName() string {
	fileNameIDMutex.Lock()
	fileNameID++
	id := fileNameID
	fileNameIDMutex.Unlock()

	return strconv.FormatInt(id, 10)
}
