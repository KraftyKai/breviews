package main

import (
	"log"
	
	"github.com/kraftykai/breview/configs"
)


func main() {
	err := configs.Init()
	
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Started!")
	
}
