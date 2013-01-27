package ServiceRegistry

// Represents a crude version numbering scheme
type Version string

// Represents all of the locations where one would find the given version of an API
type ServiceVersion struct {
	v         Version
	locations []*ServiceLocation
}

func (sr *ServiceRegistry) GetVersionFromString(vString string) (Version) {
	return Version(vString)
}

func (v1 *Version) Matches(v2 *Version) (bool) {
	if (*v1 == *v2) {
		return true
	}
	return false
}

func NewVersion(vString string) (Version) {
	return Version(vString)
}
