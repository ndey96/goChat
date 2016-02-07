package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const listenAddr = "localhost:4000"

func main() {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go match(c)
	}
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprintln(c, "Waiting for partner...")
	select {
	case partner <- c:

	case p := <-partner:
		go chat(c, p)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "You matched with someone")
	fmt.Fprintln(b, "You matched with someone")
	errc := make(chan error, 1)
	go copy(a, b, errc)
	go copy(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

func copy(a io.Writer, b io.Reader, errc chan<- error) {
	_, err := io.Copy(a, b)
	errc <- err
}
