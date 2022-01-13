package main

import (
	"fmt"
	"forum/internal/utils"
	"os"
)

func main() {
	application, err := utils.NewApp()
	if err != nil {
		fmt.Println("main:err = ", err)
		os.Exit(1)
	}
	err = application.Run()
	if err != nil {
		fmt.Println("main:err = ", err)
		os.Exit(1)
	}
}
