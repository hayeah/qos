package main

import (
	"os"

	"github.com/hayeah/qos"
)

func main() {
	qos.StartServer(os.Args[1])
}
