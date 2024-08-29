package main

import (
	"fmt"
	"go_currency_converter/Api"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <API_KEY>")
		os.Exit(1)
	}
	server := Api.NewServer(os.Args[1])
	server.Start()
}
