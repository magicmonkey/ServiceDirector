package ServiceRegistry

type ServiceVersion struct {
	major    int
	minor    int
	micro    int
	location []ServiceLocation
}
