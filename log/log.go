package mylog

import (
	"log"
	"os"
)

const (
	red   = "\033[31m"
	reset = "\033[0m"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func init() {
	Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime)
	Info = log.New(os.Stdout, "[INFO]  ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
}
