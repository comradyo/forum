package main

import (
	"forum/forum/internal/utils"
	"os"
)

func main() {
	application, err := utils.NewApp()
	if err != nil {
		os.Exit(1)
	}
	err = application.Run()
	if err != nil {
		os.Exit(1)
	}
}
