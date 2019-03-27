package geofilter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

// network,geoname_id,registered_country_geoname_id,represented_country_geoname_id,is_anonymous_proxy,is_satellite_provider
//Country, Registered Country, and Represented Country
// We now distinguish between several types of country data.
// The country is the country where the IP address is located.
// The registered_country is the country in which the IP is registered. These two may differ in some cases.
// Finally, we also include a represented_country key for some records. This is used when the IP address belongs to something like a military base.
// The represented_country is the country that the base represents. This can be useful for managing content licensing, among other uses.
func LoadCSV(ipv4Blocks io.Reader, countryCodesData io.Reader) (SliceDB, error) {
	root, err := buildCIDRRoot(ipv4Blocks)
	if err != nil {
		log.Fatal(err)
	}
	countryCodes, err := buildCountyCodes(countryCodesData)
	if err != nil {
		log.Fatal(err)
	}
	return SliceDB{root, countryCodes}, nil
}

func buildCountyCodes(r io.Reader) (map[uint32]string, error) {
	row := csv.NewReader(r)

	_, err := row.Read()
	if err != nil {
		return nil, err
	}
	out := map[uint32]string{}
	for {
		record, err := row.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) != 7 {
			fmt.Println("Expected 7 columns, got ", len(record))
			continue
		}
		geocodeID, err := strconv.ParseInt(record[0],10,32)
		if err != nil {
			log.Fatal(err)
		}
		out[uint32(geocodeID)] = record[4]
	}
	// add private networks magic code
	out[0] = "PRIVATE"
	return out, nil
}

func buildCIDRRoot(r io.Reader) (Root, error) {
	row := csv.NewReader(r)

	_, err := row.Read()
	if err != nil {
		return nil, err
	}

	var root Root
	for {
		record, err := row.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return root, err
		}
		if len(record) != 6 {
			fmt.Println("Expected 6 columns, got ", len(record))
			continue
		}
		geoID := record[2]
		if geoID == "" {
			geoID = record[1]
		}
		if record[3] != "" {
			geoID = record[3]
		}
		if geoID == "" {
			fmt.Printf("%+v; skipping\n",record)
			continue
		}
		geoIDNum, err := strconv.ParseInt(geoID, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		_, ipNet, err := net.ParseCIDR(record[0])
		if err != nil {
			log.Fatal(err)
		}
		root.Insert(ipNet, uint32(geoIDNum))

	}
	// add private networks
	privateNetworks := []string {
		"192.168.0.0/16",
		"10.0.0.0/8",
		"172.16.0.0/12",
	}
	for _,v := range privateNetworks {
		_, local, err := net.ParseCIDR(v)
		if err != nil {
			log.Fatal(err)
		}
		root.Insert(local,0)
	}
	root.Complete()
	return root, nil
}
