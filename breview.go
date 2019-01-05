package main

import (
	"fmt"
	"log"
	
	"github.com/kraftykai/breview/configs"
)


func main() {
	err := configs.Init()
	
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Let's get going!")
	fmt.Println("Port: ", configs.Values.Port)
	fmt.Println("Cfg:  ", configs.Values.File)
	fmt.Println("Host: ", configs.Values.Hostnames)
}
