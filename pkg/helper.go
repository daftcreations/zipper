package helper

import (
	"runtime"
	"strconv"
	"strings"
)

func Goid() int {
	var buf [64]byte
	id, _ := strconv.Atoi(strings.Fields(strings.TrimPrefix(string(buf[:runtime.Stack(buf[:], false)]), "goroutine "))[0])
	return id
}
