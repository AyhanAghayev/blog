package main

import (
	"anan/application"
	"fmt"
)

func main() {
	app := application.New()

	fmt.Println("Server is now listening on port :3000")
	err := app.Start()
	if err != nil {
		fmt.Println("error %w", err)
	}
}
