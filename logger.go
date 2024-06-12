package sorm

import (
  "fmt"
  "io"
  "os"
  "log"
)

func logToFile(msg string) {
  file, err := os.OpenFile("database.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

  multiWriter := io.MultiWriter(file, os.Stdout)

  log.SetOutput(multiWriter)
  log.SetFlags(log.Ldate | log.Ltime )
  log.Println(msg)
}

func logError(msg string) {
  logToFile(fmt.Sprintf("%-24v %v", "\033[31msorm.ERROR:\033[0m", msg))
}

func logInfo(msg string) {
  logToFile(fmt.Sprintf("%-15v %v", "sorm.LOG:", msg))
}

func logSuccess(msg string) {
  logToFile(fmt.Sprintf("%-24v %v", "\033[32msorm.SUCCESS:\033[0m", msg))
}
