package main

import (
	"forum/forum/internal/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	application, err := utils.NewApp()
	if err != nil {
		log.Error("main err = ", err)
		os.Exit(1)
	}
	err = application.Run()
	if err != nil {
		log.Error("main err = ", err)
		os.Exit(1)
	}
}
