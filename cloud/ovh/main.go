package main

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"

	"fmt"
)

type Continent uint8

/*
AF	Africa
AN	Antarctica
AS	Asia
EU	Europe
NA	North america
OC	Oceania
SA	South america
*/
var isoCodeToContinente = map[string]Continent{
	"AF": Africa,
	"AN": Antarctica,
	"AS": Asia,
	"EU": Europe,
	"NA": NorthAmerica,
	"OC": Oceania,
	"SA": SouthAmerica,
}

const (
//Africa ...
	Africa Continent = iota
// Antarctica ...
	Antarctica
// Asia ...
	Asia
// Europe ...
	Europe
// NorthAmerica ...
	NorthAmerica
// Oceania ...
	Oceania
// SouthAmerica ...
	SouthAmerica
//Undefined ...
	Undefined
)

//BHS1
//GRA1
//SBG1
var continentToRegion = map[Continent]string{
	Africa:       "GRA1",
	Antarctica:   "BHS1",
	Asia:         "GRA1",
	Europe:       "GRA1",
	NorthAmerica: "BHS1",
	Oceania:      "BHS1",
	SouthAmerica: "BHS1",
	Undefined:    "GRA1",
}

type ovh struct {
	opts     gophercloud.AuthOptions
	provider *gophercloud.ProviderClient
}

func (o *ovh) Authenticate() error {
	// Option 1: Pass in the values yourself
	o.opts = gophercloud.AuthOptions{
		IdentityEndpoint: "https://auth.cloud.ovh.net/v2.0",
		Username:         "mdrYmneDXkAv",
		Password:         "JsCrNmTWJ4kGgHkMpv4jdSRUKBJ6R7sX",
		TenantName:       "8484525417789870",
	}

	var err error
	o.provider, err = openstack.AuthenticatedClient(o.opts)
	if err != nil {
		return err
	}

	return nil
}

func (o *ovh) StartNode(co Continent) (ip string, err error) {
	var client *gophercloud.ServiceClient
	client, err = openstack.NewComputeV2(o.provider, gophercloud.EndpointOpts{
		Region: continentToRegion[co],
	})
	if err != nil {
		return "", err
	}


	return ip, err

}

func main() {
	ovhprovider := &ovh{}
	fmt.Println(ovhprovider.Authenticate())
}
