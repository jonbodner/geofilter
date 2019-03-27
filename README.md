# Geofilter

A simple IP-based filtering proxy, written in Go

## About
This code is a simple demonstration of using the Go standard library to build a 
performant TCP proxy server that filters out traffic based on the country of origin.

It isn't smart about traffic that has already been through another proxy, like a 
load balancer; there's nothing in the TCP/IP connection information that can be used
to figure out the original source of the request. It also doesn't know anything about
HTTP and proxy headers. It also doesn't handle IPv6 traffic.

The data comes from https://dev.maxmind.com/geoip/geoip2/geolite2/ , using the country-level
dataset. A future version might support the city-level dataset. 

## Building

Build scripts and such will come later. For now build by:

```
git clone https://github.com/jonbodner/geofilter.git
cd geofilter/filter
go build
```

## Running
To run, download the country dataset (https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country-CSV.zip),
and expand it. Then run the filter with:

```
./filter  <path_to_csv> <in_port> <out_host_port> <country_codes>
```

Where:
- `path_to_csv` is the directory that contains the `GeoLite2-Country-Blocks-IPv4.csv`
and `GeoLite2-Country-Locations-en.csv` files
- `in_port` is the port being listened on, preceded by a colon (`:`). For example, to listen on port 8000, use `:8000`.
- `out_host_port` is the hostname or IP and port of the service being proxied. For example, to proxy ssh on server running
the filer, use `localhost:22`
- `country_codes` is a comma-separated list of two-letter ISO 3166-1 alpha-2 country codes that are allowed to pass the filter. The special
country code `PRIVATE` allows traffic from private networks (192.168.0.0/24, 10.0.0.0/8, and 172.16.0.0/12)

So, to allow access from private networks and the US, listening on port 8000, redirecting to SSH on the local machine, and using
CSV files in the current directory, you would use the command line:

```
./filter . :8000 localhost:22 US,PRIVATE
```



