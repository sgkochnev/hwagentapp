package logger

import (
	"log"
	"os"
)

func Log(v ...any) {
	f, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println(v...)
}
