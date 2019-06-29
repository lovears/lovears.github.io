package logger

import (
	"log"
	"os"
	"time"
	"fmt"
	"sync"
	"github.com/robfig/cron"
)

var one = sync.Once{}
var logger *log.Logger
var std *log.Logger
var debug = true

func init() {
	fileName := "./" + time.Now().Format("20060102") + ".log"
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	one.Do(func() {
		logger = log.New(logFile, "crucian_blog_", log.LstdFlags|log.Lshortfile)
		std = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
		c := cron.New()
		c.AddFunc("@daily", schedule)
	})
}

func schedule() {
	fileName := "./" + time.Now().Format("20060102") + ".log"
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(logFile)
}

func Info(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.Output(2, s)
	if debug {
		std.Output(2, s)
	}
}

func Panic(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.Output(2, s)
	if debug {
		std.Output(2, s)
	}
	panic(s)
}
