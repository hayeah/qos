package main

import (
	"log"
	"os"

	"github.com/hayeah/qos"
)

func main() {
	addr := os.Args[1]

	client, err := qos.NewClient(addr)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Get(10, os.Stdout)
	if err != nil {
		log.Println(err)
	}
	err = client.Get(50, os.Stdout)
	if err != nil {
		log.Println(err)
	}

}
