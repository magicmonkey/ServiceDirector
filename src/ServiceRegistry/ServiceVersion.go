package ServiceRegistry

type Version struct {
	Major    int
	Minor    int
	Micro    int
}

type ServiceInstance struct {
	Version  Version
	Location ServiceLocation
}

type ServiceVersion struct {
	v         Version
	locations []ServiceLocation
}
