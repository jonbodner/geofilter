package geofilter

import "net"

type DB interface {
	Code(ip net.IP) (string, error)
}
