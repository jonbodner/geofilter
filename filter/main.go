package main

import (
	"fmt"
	"github.com/jonbodner/geofilter"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "usage: %s <path_to_csv> <in_port> <out_host_port> <country_codes>\n", os.Args[0])
		os.Exit(1)
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	db, err := func(fileName string) (geofilter.DB, error) {
		blocksFile, err := os.Open(fileName + "/GeoLite2-Country-Blocks-IPv4.csv")
		if err != nil {
			return nil, err
		}
		defer blocksFile.Close()
		ccFile, err := os.Open(fileName + "/GeoLite2-Country-Locations-en.csv")
		if err != nil {
			return nil, err
		}
		defer ccFile.Close()
		db, err := geofilter.LoadCSV(blocksFile, ccFile)
		if err != nil {
			return nil, err
		}
		return db, nil
	}(os.Args[1])
	if err != nil {
		panic(err)
	}
	codes := map[string]bool{}
	for _, v := range strings.Split(os.Args[4], ",") {
		codes[v] = true
	}
	geofilter.ListenAndProcess(os.Args[2], os.Args[3], func(ip net.IP) bool {
		if ip.IsLoopback() {
			return true
		}
		code, err := db.Code(ip)
		if err != nil {
			fmt.Println(err)
			return false
		}
		fmt.Println(code)
		return codes[code]
	})
}
