// This package tracks the data structures for the entire Service Registry, including the contained Services and
// their locations
package ServiceRegistry

import (
	"net/url"
	"math/rand"
)

// An API, of which there may be many versions, each of which may have many locations across which you want to balance
// the load
type Service struct {
	Name     string
	Versions []*ServiceVersion
}

// The top-level object which contains all of the services available
type ServiceRegistry struct {
	Services []*Service
}

// An abstraction of the location of a service; currently only allows a URL to be used
type ServiceLocation struct {
	Location url.URL
}

func (sr *ServiceRegistry) GetServiceWithName(name string, createIfNotExist bool) (service *Service, created bool) {
	for i, value := range sr.Services {
		if (value.Name == name) {
			return sr.Services[i], false
		}
	}

	if createIfNotExist {
		// Need to make a new Service
		s := NewService(name)
		sr.Services = append(sr.Services, s)
		return s, true
	}

	return nil, false
}

func NewService(name string) (*Service) {
	s := new(Service)
	s.Name = name
	return s
}

func NewLocation(u *url.URL) (*ServiceLocation) {
	sl := ServiceLocation{*u}
	return &sl
}

func (s *Service) GetLocationsForVersion(v Version) ([]*ServiceLocation) {
	for i, value := range s.Versions {
		if (value.Version.Matches(&v, false)) {
			return s.Versions[i].locations
		}
	}
	return nil
}

func (s *Service) GetLocationForVersion(v Version) (*ServiceLocation) {
	sls := s.GetLocationsForVersion(v)
	if sls == nil {
		return nil
	}
	// Randomly pick a version
	vNum := rand.Intn(len(sls))
	return sls[vNum]
}

func (s *Service) getVersion(v Version) (*ServiceVersion) {
	for i, value := range s.Versions {
		if (value.Version.Matches(&v, true)) {
			return s.Versions[i]
		}
	}

	// Need to make a new ServiceVersion
	sv := ServiceVersion{v, nil}
	s.Versions = append(s.Versions, &sv)
	return &sv
}

func (s *Service) AddServiceInstance(v Version, sl *ServiceLocation) {
	sv := s.getVersion(v)
	sv.locations = append(sv.locations, sl)
}
