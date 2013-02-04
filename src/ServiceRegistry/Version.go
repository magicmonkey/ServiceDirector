package ServiceRegistry

import (
	"strings"
)

// Represents a crude version numbering scheme
type Version string

// Represents all of the locations where one would find the given version of an API
type ServiceVersion struct {
	Version   Version
	Locations []*ServiceLocation
}

func (sr *ServiceRegistry) GetVersionFromString(vString string) (Version) {
	return Version(vString)
}

func (v1 *Version) Matches(v2 *Version, strict bool) (bool) {

	// TODO: Sort in order

	// Strict matching is easy; just compare the strings
	if strict {
		if (*v1 == *v2) {
			return true
		}
		return false
	}

	// Non-strict matching needs the parts checking individually
	v1Parts := strings.Split(string(*v1), ".")
	v2Parts := strings.Split(string(*v2), ".")

	// Check for an impossible-to-satisfy case
	if len(v2Parts) > len(v1Parts) {
		return false
	}

	for i, _ := range v1Parts {
		if len(v2Parts) > i {
			if v1Parts[i] == v2Parts[i] {
				continue
			}
			return false
		}
		return true
	}
	return true
}

func NewVersion(vString string) (Version) {
	return Version(vString)
}
