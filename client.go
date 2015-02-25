package qos

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func StartClient(port string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	startTime := time.Now()

	log.Println("start receiving bytes")
	buf := make([]byte, 1024*4)

	bytesRead := 0
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		} else {
			bytesRead += n
		}

	}

	fmt.Printf("total bytes read: %v\n", bytesRead)
	d := time.Since(startTime)
	fmt.Printf("time taken: %v\n", d.Seconds())
	fmt.Printf("kb/s: %v\n", float64(bytesRead/KB)/d.Seconds())
	return nil
}
