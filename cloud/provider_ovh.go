package cloud

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

type ovh struct {
	opts     gophercloud.AuthOptions
	provider *gophercloud.ProviderClient
}

func (o *ovh) Authenticate() error {
	// Option 1: Pass in the values yourself
	o.opts = gophercloud.AuthOptions{
		IdentityEndpoint: "https://auth.cloud.ovh.net/v2.0",
		Username:         "mdrYmneDXkAv",
		Password:         "6XEUqjbMkfyHqXgy8nYbErQvsYcFxqfr",
		TenantID:         "8484525417789870",
	}
	var err error
	o.provider, err = openstack.AuthenticatedClient(o.opts)
	if err != nil {
		return err

	}

	return nil
}
