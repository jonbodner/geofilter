package geofilter

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Returns true to indicate the IP is allowed, false if the IP is blocked
type FilterFunc func(ip net.IP) bool

func ListenAndProcess(inPort string, outHostPort string, allow FilterFunc) {
	ln, err := net.Listen("tcp", inPort)
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
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			//handle error
			fmt.Println(err)
			continue
		}
		ip := net.ParseIP(host)
		if allow(ip) {
			go func(conn net.Conn) {
				defer conn.Close()
				handleConnection(conn, outHostPort)
			}(conn)
		} else {
			err = conn.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func handleConnection(inConn net.Conn, outHostPort string) {
	outConn, err := net.Dial("tcp", outHostPort)
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}
	defer outConn.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := io.Copy(inConn, outConn)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(outConn, inConn)
		if err != nil {
			fmt.Println(err)
		}
	}()
	wg.Wait()
}
