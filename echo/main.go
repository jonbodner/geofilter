package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func Server(port string) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println(err)
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			_, err := io.Copy(conn, conn)
			if err != nil {
				fmt.Println(err)
			}
		}(conn)
	}
}

func Client(hostPort string) {
	outConn, err := net.Dial("tcp", hostPort)
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}
	defer outConn.Close()
	for i := 1024; i <= 1024*1024; i = i * 2 {
		b := make([]byte,i)
		for k := 0;k<i;k++ {
			b[k] = byte(k % 256)
		}
		out := make([]byte,i)
		start := time.Now()
		for j := 0; j < 10; j++ {
			n, err := outConn.Write(b)
			if err != nil {
				panic(err)
			}
			if n != i {
				panic("should write all bytes")
			}
			total := 0
			for total != i {
				n, err := outConn.Read(out)
				if err != nil {
					panic(err)
				}
				total += n
			}
		}
		end := time.Now()
		fmt.Printf("average time to write %d bytes is %v\n", i, end.Sub(start)/10)
	}
}

func main() {
	if len(os.Args) == 3 {
		if os.Args[1] == "server" {
			Server(os.Args[2])
		}
		if os.Args[1] == "client" {
			Client(os.Args[2])
		}
		return
	}
	fmt.Fprintln(os.Stderr, "usage: "+os.Args[0]+" server|client port")
}
