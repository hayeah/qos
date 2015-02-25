package main

import (
	"fmt"
	"log"

	"os"

	"github.com/hayeah/qos"
)

func main() {
	port := os.Args[1]
	fmt.Printf("dial: %v\n", port)

	err := qos.StartClient(port)
	if err != nil {
		log.Println(err)
	}
}
