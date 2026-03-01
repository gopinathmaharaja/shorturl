package logger

import (
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	WarnLog  *log.Logger
	ErrorLog *log.Logger
)

func Init() {
	InfoLog = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)
	WarnLog = log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Lshortfile)
	ErrorLog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
}
