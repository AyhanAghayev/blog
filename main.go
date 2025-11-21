package main

import (
	"anan/application"
	"fmt"
)

func main() {
	app := application.New()

	err := app.Start()
	if err != nil {
		fmt.Println("error %w", err)
	}
}
