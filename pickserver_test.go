package main

import (
	"log"
	"testing"
)

var serverIPs = []string{
	"151.237.84.46", //usa
	"212.13.49.186", //pt
	"84.124.11.207", //es
	"85.13.13.254", //au
}

var clientIP = "5.141.9.86"

func TestNearIp(t *testing.T) {
	serversCoords := make(map[string]Coordinates)

	for _, ip := range serverIPs {
		serversCoords[ip] = IPToCoords(ip)
	}

	var distanceIps = make(map[float64]string)

	ipLoc := IPToCoords(clientIP)

	//calculate userip distance to each server
	for ip := range serversCoords {
		dst := Distance(serversCoords[ip].Latitude, ipLoc.Longitude, ipLoc.Latitude, serversCoords[ip].Longitude)
		distanceIps[dst] = ip
	}
	//compare distance to each nodes

	//get the nearest server
	log.Println(nearIP(distanceIps))

}
