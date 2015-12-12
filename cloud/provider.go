package cloud

//Provider represents an interface to  Cloud Providers (DO/OVH)
type Provider interface {
	Authenticate() error
	//StartNode starts a node on a given Continent and returns the nodeIP
	StartNode(Continent) (ip string, err error)
	StopNode(ip string) error
}
