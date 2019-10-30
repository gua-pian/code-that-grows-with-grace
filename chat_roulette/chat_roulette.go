package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const listenAddr = ":4000"

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
	fmt.Fprintln(c, "Waiting for a parnter...")
	select {
	case partner <- c:
	// now handled by other goroutine
	case p := <-partner:
		chat(p, c)
	}
}
func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

func cp(w, r io.ReadWriteCloser, errc chan error) {
	_, err := io.Copy(w, r)
	errc <- err
}
