package qos

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

func StartServer(port string) error {
	so, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	for {
		conn, err := so.Accept()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("new connection")
			handleServerLoad(conn)
		}
	}
}

const (
	KB = 1024
	MB = 1024 * 1024
)

var randKB []byte

func init() {
	randKB = make([]byte, 1024)
	for i := 0; i < len(randKB); i++ {
		randKB[i] = byte(97 + (rand.Int() % 26))
	}
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func handleServerLoad(conn net.Conn) error {
	defer conn.Close()
	var err error
	// var buf []byte

	nbytes := 0
	// // write 10M
	for i := 0; i < 1024; i++ {
		n, err := conn.Write(randKB)
		if err != nil {
			log.Println(err)
			break
		}
		nbytes += n
	}
	fmt.Printf("total bytes written: %v\n", nbytes)

	// n, err := io.Copy(os.Stdout, conn)
	// fmt.Printf("read:%v\n", n)
	// if err != nil {
	// 	log.Println(err)
	// }

	// for {
	// 	var buf []byte
	// 	conn.Rea
	// 	_, err := conn.Read(buf)
	// 	fmt.Printf("read:%v\n", buf)
	// 	if err != nil {
	// 		log.Println(err)
	// 		break
	// 	}
	// }

	return err
	// read 1M
}
