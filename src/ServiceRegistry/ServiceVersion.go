package ServiceRegistry

type Version struct {
	major    int
	minor    int
	micro    int
}

type ServiceInstance struct {
	Version  Version
	Location ServiceLocation
}

type ServiceVersion struct {
	v         Version
	locations []ServiceLocation
}
