package qos

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

var (
	CmdBytes = "BYTES"
	CmdGet   = "GET"
)

func StartServer(port string) error {
	addr := fmt.Sprintf(":%v", port)
	so, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Println("listening:", addr)
	for {
		conn, err := so.Accept()

		if err != nil {
			log.Println(err)
		} else {
			log.Println("client connected")
			server := &Server{conn: conn}
			go server.commandLoop()
		}
	}
}

const (
	KB = 1024
	MB = 1024 * 1024
)

var randChunk []byte

type Server struct {
	conn net.Conn
}

// type ServerHandler {
// }

var (
	UnknownCommandError = []byte("ERROR")
	ClientError         = []byte("CLIENT_ERROR")
)

func (s *Server) commandLoop() {
	// what's the client-side close error?
	// peek next line
	conn := s.conn
	defer conn.Close()

	lines := bufio.NewScanner(conn)

	for {

		if ok := lines.Scan(); !ok {
			if err := lines.Err(); err != nil {
				log.Println(err)
			} else {
				log.Println("client disconnected")
			}
			break
		}

		line := lines.Text()

		parts := strings.Split(line, " ")

		cmd := strings.ToUpper(parts[0])
		args := parts[1:]

		switch cmd {
		case "GET":
			// bytes <nbytes>
			// VALUE <nbytes>\r\n
			//
			if len(args) != 1 {
				s.replyClientError("Invalid number of arguments.")
				continue
			}

			nbytes, err := strconv.Atoi(args[0])
			if err != nil {
				s.replyClientError("Number of bytes not an integer.")
				continue
			}

			log.Printf("writing %v bytes to client\n", nbytes)
			// not quit sure what to do if there's error when outputting bytes. close connection?
			err = s.replyBytes(nbytes)
			if err != nil {
				log.Println(err)
			}

			// for n, nwritten := 0; n < nbytes; n + nwritten {
			// 	bytes = randKB
			// 	conn.Write(randKB)
			// }

		default:
			// "ERROR\r\n"
			conn.Write(UnknownCommandError)
			conn.Write([]byte("\r\n"))
		}
	}

	return
}

func (s *Server) replyBytes(nbytes int) error {
	conn := s.conn

	var err error
	conn.Write([]byte("BYTES "))
	conn.Write([]byte(strconv.Itoa(nbytes)))
	conn.Write([]byte("\r\n"))

	nwritten := 0
	chunkSize := len(randChunk)
	for nwritten < nbytes {
		bytesRemain := nbytes - nwritten
		chunk := randChunk
		if bytesRemain < chunkSize {
			chunk = chunk[0:bytesRemain]
		}

		n, err := conn.Write(chunk)
		if err != nil {
			break
		}

		nwritten += n
	}

	conn.Write([]byte("\r\n"))

	return err
}

func (s *Server) replyClientError(msg string) {
	conn := s.conn
	conn.Write(ClientError)
	conn.Write([]byte(" "))
	conn.Write([]byte(msg))
	conn.Write([]byte("\r\n"))
}

func (s *Server) handlePing() {

}

func init() {
	randChunk = make([]byte, 1024*4)
	for i := 0; i < len(randChunk); i++ {
		randChunk[i] = byte(97 + (rand.Int() % 26))
	}
	log.SetFlags(log.Lshortfile | log.Ltime)
}
