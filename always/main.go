package main

import (
	"fmt"
	"github.com/jonbodner/geofilter"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <host and port to proxy>\n",os.Args[0])
		os.Exit(1)
	}
	geofilter.ListenAndProcess(":8000", os.Args[1], func(ip net.IP) bool {
		return true
	})
}
