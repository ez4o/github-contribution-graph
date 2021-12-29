package main

import (
	"encoding/json"
	"fmt"
	"os"
	"server"
	"server/model"
)

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error while opening config:", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := model.Configuration{}

	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Error while parsing config:", err)
	}

	s := &server.Server{GitHubToken: configuration.GitHubToken}
	s.Start()
}
