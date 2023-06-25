package logger

import (
	"fmt"
	"time"
)

func Info(line string) {
	log(line, "INFO")
}

func Warn(line string) {
	log(line, "WARN")
}

func log(line, level string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%v [%v] %v\n", now, level, line)
}
