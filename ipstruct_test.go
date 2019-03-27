package geofilter

import (
	"errors"
	"net"
	"testing"
)

func TestRoot(t *testing.T) {
	data := []struct {
		cidr string
		code uint32
	} {
		{"1.0.0.0/24", 1},
		{"1.0.1.0/24", 2},
		{"1.0.2.0/23", 3},
		{"1.0.4.0/22", 4},
		{"1.0.8.0/21", 5},
		{"1.0.16.0/20", 6},
	}
	var r Root
	for _, v := range data {
		_, ipnet, err := net.ParseCIDR(v.cidr)
		if err != nil {
			t.Fatal(err)
		}
		r.Insert(ipnet, v.code)
	}
	r.Complete()
	testData := []struct {
		ip string
		code uint32
		err error
	} {
		{"1.0.0.0", 1,nil},
		{"1.0.0.128", 1,nil},
		{"1.0.1.0", 2,nil},
		{"1.0.1.128", 2,nil},
		{"1.0.2.0", 3,nil},
		{"1.0.3.128", 3,nil},
		{"1.0.3.0", 3,nil},
		{"1.0.3.128", 3,nil},
		{"1.0.4.0", 4,nil},
		{"1.0.5.0", 4,nil},
		{"1.0.6.0", 4,nil},
		{"1.0.7.0", 4,nil},
		{"1.0.8.0", 5,nil},
		{"1.0.9.0", 5,nil},
		{"1.0.10.0", 5,nil},
		{"1.0.11.0", 5,nil},
		{"1.0.12.0", 5,nil},
		{"1.0.13.0", 5,nil},
		{"1.0.14.0", 5,nil},
		{"1.0.15.0", 5,nil},
		{"1.0.16.0", 6,nil},
		{"1.0.20.0", 6,nil},
		{"1.0.30.0", 6,nil},
		{"1.0.32.0", 0,errors.New("No country for ip 1.0.32.0")},
	}
	for _, v := range testData {
		ip := net.ParseIP(v.ip)
		code, err := r.Code(ip)
		if err != nil && v.err == nil{
			t.Error(err)
		}
		if err == nil && v.err != nil {
			t.Errorf("Expected error %s, got nil",v.err.Error())
		}
		if err != nil && v.err != nil && err.Error() != v.err.Error() {
			t.Errorf("Expected error %s, got %s", v.err.Error(), err.Error())
		}
		if code != v.code {
			t.Errorf("Expected %d, got %d for %s",v.code, code, v.ip)
		}
	}
}
