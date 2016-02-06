package cloud

import (
	"strconv"
	"testing"
	"time"

	"github.com/digitalocean/godo"
)

func TestDigitalOceanAuthenticate(t *testing.T) {
	do := &digitalOcean{}
	err := do.Authenticate()
	if err != nil {
		t.Errorf("failed with error %s", err)
		t.FailNow()
	}
    opt := &godo.ListOptions{}
     panic( do.c.Sizes.List(opt))
    
  
}

func TestDigitalOceanUpdateRegions(t *testing.T) {
	do := &digitalOcean{}
	err := do.Authenticate()
	if err != nil {
		t.Errorf("failed to auth with error %s", err)
		t.FailNow()
	}

	err = do.updateAvailableRegions()
	if err != nil {
		t.Errorf("failed to update available regions with error %s", err)
		t.FailNow()
	}

}

func TestDigitalOceanNodeCreateAndDelete(t *testing.T) {
	do := &digitalOcean{}
	err := do.Authenticate()
	if err != nil {
		t.Errorf("failed to auth with error %s", err)
		t.FailNow()
	}

	dropletName := "edgefrffh"

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: "fra1",
		Size:   "512mb",
		Image: godo.DropletCreateImage{

			ID: 14575753,
		},
	}
	newDroplet, _, err := do.c.Droplets.Create(createRequest)
	if err != nil {
		t.Errorf("Something bad happened: %s\n\n", err)
		t.FailNow()

	}

	time.Sleep(time.Second * 3)

	err = do.StopNode(strconv.Itoa(newDroplet.ID))
	if err != nil {
		t.Errorf("Something bad happened: %s\n\n", err)
		t.FailNow()

	}

}
