package main

import (
	"encoding/json"
	"log"
	"os"
	"server"
	"server/model"
)

func init() {
	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}

	log.SetOutput(logFile)
	log.SetPrefix("[log] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Error while opening config:", err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	configuration := model.Configuration{}

	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatalln("Error while parsing config:", err)
	}

	s := &server.Server{GitHubToken: configuration.GitHubToken}
	s.Start()
}
