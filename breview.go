package main

import (
	"log"
	
	"github.com/kraftykai/breview/configs"
	"github.com/kraftykai/breview/server"
)


func main() {
	
	cfg, err := configs.Init()
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Init(cfg); err != nil {
		log.Fatal(err)
	}
	log.Printf("Started!")
	
}
