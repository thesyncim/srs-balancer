package cloud

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"labix.org/v2/mgo/bson"
)

var (
	pat = "9968278d444269ac52a80a0bd0c72f96ce37feb7b948d172b763ca170b92994c"
	tokenSrc = &tokenSource{
		AccessToken: pat,
	}
	doServerSize = "512mb"
)

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type digitalOcean struct {
	maxNodes         int
	c                *godo.Client
	AvailableRegions map[string]Continent
}

/*
nyc1", Name:"New York 1"
ams1", Name:"Amsterdam 1
sfo1", Name:"San Francis""
nyc2", Name:"New York 2"
ams2", Name:"Amsterdam 2
sgp1", Name:"Singapore 1
lon1", Name:"London 1"
nyc3", Name:"New York 3"
ams3", Name:"Amsterdam 3
fra1", Name:"Frankfurt 1
tor1", Name:"Toronto 1"
*/

var allregions = map[string]Continent{
	"nyc1": NorthAmerica,
	"nyc3": NorthAmerica,
	"tor1": NorthAmerica,
	"sfo1": NorthAmerica,
	"nyc2": NorthAmerica,
	"ams1": Europe,
	"ams2": Europe,
	"lon1": Europe,
	"ams3": Europe,
	"fra1": Europe,
	"sgp1": Asia,
}

func (do *digitalOcean) updateAvailableRegions() error {
	lopts := &godo.ListOptions{}
	lopts.PerPage = 100

	//assume all regions are available
	do.AvailableRegions = allregions

	//pick all regions
	r, _, err := do.c.Regions.List(lopts)
	if err != nil {
		return err
	}

	//delete unavailable regions if offline or doesnt contain desired machine size
	for i := range r {
		if !r[i].Available || !contains(r[i].Sizes, doServerSize) {
			delete(do.AvailableRegions, r[i].Slug)
		}

	}

	return nil
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func (do *digitalOcean) pickDCRegion(c Continent) (string, error) {
	err := do.updateAvailableRegions()
	if err != nil {
		return "", err
	}
	var regions []string
	//get available regions for a datacenter
	for i := range do.AvailableRegions {
		if do.AvailableRegions[i] == c {
			regions = append(regions, i)
		}
	}

	if len(regions) == 0 {
		return "", ErrUnavailableDatacenter
	}
	//pick one random
	r := random(0, len(regions) - 1)
	return regions[r], nil
}

func (do *digitalOcean) newDropletName() string {
	return "digitalOcean-" + bson.NewObjectId().Hex()

}

func (do *digitalOcean) StartNode(c Continent) (string, error) {

	region, err := do.pickDCRegion(c)
	if err != nil {
		return "", err
	}
	createRequest := &godo.DropletCreateRequest{
		Name:   do.newDropletName(),
		Region: region,
		Size:   doServerSize,
		Image: godo.DropletCreateImage{
			ID: 14575753,
		},
	}
	newDroplet, _, err := do.c.Droplets.Create(createRequest)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(newDroplet.ID), nil
}

func (do *digitalOcean) StopNode(id string) error {
	intid, _ := strconv.Atoi(id)
	_, err := do.c.Droplets.Delete(intid)

	return err
}

func (do *digitalOcean) Authenticate() error {
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSrc)
	do.c = godo.NewClient(oauthClient)

	acc, _, err := do.c.Account.Get()
	if err != nil {
		return err
	}
	do.maxNodes = acc.DropletLimit
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

var _ = Provider(&digitalOcean{})
