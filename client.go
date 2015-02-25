package qos

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
)

var (
	ErrBadServerResponse = errors.New("Invalid response from server")
)

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{bufr: bufio.NewReader(conn), conn: conn}, nil
	// startTime := time.Now()

	// log.Println("start receiving bytes")
	// buf := make([]byte, 1024*4)

	// bytesRead := 0
	// for {
	// 	n, err := conn.Read(buf)
	// 	if err != nil {
	// 		if err != io.EOF {
	// 			log.Println(err)
	// 		}
	// 		break
	// 	} else {
	// 		bytesRead += n
	// 	}
	// }

	// fmt.Printf("total bytes read: %v\n", bytesRead)
	// d := time.Since(startTime)
	// fmt.Printf("time taken: %v\n", d.Seconds())
	// fmt.Printf("kb/s: %v\n", float64(bytesRead/KB)/d.Seconds())
}

type Client struct {
	conn net.Conn
	bufr *bufio.Reader
}

func (c *Client) Get(nbytes int, w io.Writer) error {
	// nbytes, err := io.CopyN(w, c.conn, n)
	c.request(CmdGet, strconv.Itoa(nbytes))
	cmd, args, err := c.receiveLine()
	if err != nil {
		return err
	}

	if cmd != CmdBytes {
		return ErrBadServerResponse
	}

	nbytesToReceive, err := strconv.Atoi(args[0])
	if err != nil {
		return ErrBadServerResponse
	}

	err = c.receiveBytes(nbytesToReceive, w)

	// consume \r\n
	c.bufr.ReadByte()
	c.bufr.ReadByte()

	return err
}

func (c *Client) receiveBytes(nbytes int, w io.Writer) error {
	_, err := io.CopyN(w, c.bufr, int64(nbytes))
	return err
}

func (c *Client) receiveLine() (cmd string, args []string, err error) {
	var lines []string
	for {
		lineFrag, isPrefix, err := c.bufr.ReadLine()
		if err != nil {
			return "", nil, err
		}

		lines = append(lines, string(lineFrag))

		if !isPrefix {
			break
		}
	}

	line := strings.Join(lines, "")

	return parseLine(line)
}

func (c *Client) request(cmd string, args ...string) error {
	conn := c.conn
	var err error
	_, err = conn.Write([]byte(cmd))
	for _, arg := range args {
		_, err = conn.Write([]byte(" "))
		_, err = conn.Write([]byte(arg))
		if err != nil {
			break
		}
	}
	_, err = conn.Write([]byte("\r\n"))
	return err
}
