package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/hayeah/qos"
)

type NullWriter struct{}

func (n *NullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func NewCounterWriter(w io.Writer) *CounterWriter {
	return &CounterWriter{w: w}
}

type CounterWriter struct {
	BytesWritten int
	w            io.Writer
}

func (c *CounterWriter) Reset() {
	c.BytesWritten = 0
}

func (c *CounterWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.BytesWritten += n
	return n, err
}

// sample write every n milliseconds
func benchWrite(sampleEvery int, w io.Writer, fn func(cw io.Writer) error) error {
	start := time.Now()
	cw := &CounterWriter{w: w}
	done := make(chan error)
	go func() {
		err := fn(cw)
		done <- err
	}()

	ticker := time.NewTicker(time.Duration(sampleEvery) * time.Millisecond)
	defer ticker.Stop()

	var err error
	nbytes := 0
	lastTick := start
loop:
	for {
		select {
		case t := <-ticker.C:
			// nbytes = cw.BytesWritten - nbytes
			dt := t.Sub(lastTick)

			bytesWritten := cw.BytesWritten - nbytes
			nbytes = cw.BytesWritten
			lastTick = t

			log.Println("bytes:", bytesWritten)
			log.Printf("rate(kb/s): %.2f\n", float64(bytesWritten/1024)/dt.Seconds())
		case err = <-done:
			break loop
		}
	}

	// report throughput
	nbytes = cw.BytesWritten
	log.Println("total bytes:", cw.BytesWritten)
	log.Printf("throughput(kb/s): %.2f\n", float64(nbytes/1024)/time.Since(start).Seconds())
	return err
}

func main() {
	addr := os.Args[1]

	client, err := qos.NewClient(addr)

	if err != nil {
		log.Fatal(err)
	}

	// 1G
	err = benchWrite(100, &NullWriter{}, func(w io.Writer) error {
		return client.Get(1024*1024*1024, w)
	})

	if err != nil {
		log.Println(err)
	}

}
