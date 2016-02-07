package main

import (
    "io"
    "log"
    "net"
    "fmt"
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
    fmt.Fprintln(c, "Waiting for partner")
    select {
        case partner <- c:

        case p := <-partner:
            chat(c, p)
    }
}

func chat(a,b io.ReadWriteCloser) {
    fmt.Fprintln(a, "You matched with someone")
    fmt.Fprintln(b, "You matched with someone")
    go io.Copy(a,b)
    go io.Copy(b,a)
}