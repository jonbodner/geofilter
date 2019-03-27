package geofilter

import (
	"errors"
	"net"
	"sort"
)

type Cell struct {
	val  uint32
	mask *net.IPNet
	code uint32
}

type Root []Cell

func (r *Root) Insert(mask *net.IPNet, code uint32) {
	*r = append(*r, Cell{val: ip2uint32(mask.IP), mask: mask, code: code})
}

func (r *Root) Complete() {
	sort.Slice(*r, func(i, j int) bool {
		return (*r)[i].val > (*r)[j].val
	})
}

func (r Root) Code(ip net.IP) (uint32, error) {
	ipv4 := ip.To4()
	if ipv4 == nil {
		return 0, errors.New("Not an ip in IPv4: " + ip.String())
	}
	val := ip2uint32(ipv4)
	pos := sort.Search(len(r), func(i int) bool {
		return r[i].val <= val
	})
	if r[pos].mask.Contains(ip) {
		return r[pos].code, nil
	}
	return 0, errors.New("No country for ip " + ip.String())
}

func ip2uint32(ip net.IP) uint32 {
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

type SliceDB struct {
	root         Root
	countryCodes map[uint32]string
}

func (db SliceDB) Code(ip net.IP) (string, error) {
	geonameID, err := db.root.Code(ip)
	if err != nil {
		return "", err
	}
	return db.countryCodes[geonameID], nil
}
