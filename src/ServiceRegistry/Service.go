// This package tracks the data structures for the entire Service Registry, including the contained Services and
// their locations
package ServiceRegistry

import (
	"math/rand"
	"log"
)

// An API, of which there may be many versions, each of which may have many locations across which you want to balance
// the load
type Service struct {
	Name     string
	Versions []*ServiceVersion
	serviceRegistry *ServiceRegistry
}

// The top-level object which contains all of the services available
type ServiceRegistry struct {
	Services   []*Service
	Name string
	serviceRegistryUpdateChan []chan ServiceRegistry
}

// An abstraction of the location of a service; currently only allows a URL to be used
type ServiceLocation struct {
	Location string
}

func NewServiceRegistry(name string, sru chan ServiceRegistry) (*ServiceRegistry) {
	sr := ServiceRegistry{}
	sr.Name = name
	sr.RegisterUpdateChannel(sru)
	sr.SendRegistryUpdate()
	return &sr
}

func NewLocation(u string) (*ServiceLocation) {
	sl := ServiceLocation{u}
	return &sl
}

func (sr *ServiceRegistry) SendRegistryUpdate() {
	log.Printf("[Service] Send update to all concerned (from %v)\n", sr)
	for _, c := range sr.serviceRegistryUpdateChan {
		log.Println("[Service] Sending update...")
		c <- *sr
	}
}

func (sr *ServiceRegistry) GetServiceWithName(name string, createIfNotExist bool) (service *Service, created bool) {
	for i, value := range sr.Services {
		if (value.Name == name) {
			return sr.Services[i], false
		}
	}

	if createIfNotExist {
		// Need to make a new Service
		s := sr.NewService(name)
		return s, true
	}

	return nil, false
}

func (sr *ServiceRegistry) NewService(name string) (*Service) {
	s := new(Service)
	s.Name = name
	s.serviceRegistry = sr
	sr.Services = append(sr.Services, s)
	sr.SendRegistryUpdate()
	return s
}

func (sr *ServiceRegistry) RegisterUpdateChannel(sru chan ServiceRegistry) {
	log.Println("[Service] Adding update channel...")
	sr.serviceRegistryUpdateChan = append(sr.serviceRegistryUpdateChan, sru)
}

func (s *Service) GetLocationsForVersion(v Version) ([]*ServiceLocation) {
	for i, value := range s.Versions {
		if (value.Version.Matches(&v, false)) {
			return s.Versions[i].Locations
		}
	}
	return nil
}

func (s *Service) GetLocationsForVersionString(v string) ([]*ServiceLocation) {
	return s.GetLocationsForVersion(s.serviceRegistry.getVersionFromString(v))
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

func (s *Service) GetLocationForVersionString(v string) (*ServiceLocation) {
	return s.GetLocationForVersion(s.serviceRegistry.getVersionFromString(v))
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
	sv.Locations = append(sv.Locations, sl)
	s.serviceRegistry.SendRegistryUpdate()
}

func (s *Service) SetServiceRegistry(sr *ServiceRegistry) {
	s.serviceRegistry = sr
}
