package cloud

import (

	"net"
	"path/filepath"

	"github.com/thesyncim/geoip2-golang"
)

var db *geoip2.Reader

func init() {
	var err error
	db, err = geoip2.Open(filepath.Join("data", "GeoLite2-City.mmdb"))

	if err != nil {
		log.Fatal(err)
	}
}

//Coordinates ... represents
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

//IPToCoords returns coordinates based on input IP
func IPToCoords(ipstr string) (Coordinates, Continent) {

	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return Coordinates{}, Undefined
	}
	record, err := db.City(ip)
	if err != nil {
		log.Println(err)
		return Coordinates{}, Undefined
	}

	return Coordinates{
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
	}, isoCodeToContinente[record.Continent.Code]
}
