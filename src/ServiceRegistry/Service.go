package ServiceRegistry

import (
	"net/url"
	"math/rand"
	"strings"
	"strconv"
)

type Version struct {
	Major    int64
	Minor    int64
	Micro    int64
}

type ServiceInstance struct {
	Version  Version
	Location *ServiceLocation
}

type ServiceVersion struct {
	v         Version
	locations []*ServiceLocation
}

type Service struct {
	Name     string
	Versions []*ServiceVersion
}

type ServiceRegistry struct {
	Services []*Service
}

type ServiceLocation struct {
	Location url.URL
}

func (sr *ServiceRegistry) GetVersionFromString(vString string) (Version) {
	versionParts := strings.Split(vString, ".")
	vp0, err := strconv.ParseInt(versionParts[0], 10, 8)
	if err != nil {
		panic(err)
	}
	vp1, err := strconv.ParseInt(versionParts[1], 10, 8)
	if err != nil {
		panic(err)
	}
	vp2, err := strconv.ParseInt(versionParts[2], 10, 8)
	if err != nil {
		panic(err)
	}
	return Version{vp0, vp1, vp2}
}

func (v1 *Version) Matches(v2 *Version) (bool) {
	if (v1.Major == v2.Major && v1.Minor == v2.Minor && v1.Micro == v2.Micro) {
		return true
	}
	return false
}

func (sr *ServiceRegistry) GetServiceWithName(name string) (*Service) {
	for i, value := range sr.Services {
		if (value.Name == name) {
			return sr.Services[i]
		}
	}

	// Need to make a new Service
	s := NewService(name)
	sr.Services = append(sr.Services, s)
	return s
}

func NewService(name string) (*Service) {
	s := new(Service)
	s.Name = name
	return s
}

func NewVersion(major int64, minor int64, micro int64) (Version) {
	v := Version{major, minor, micro}
	return v
}

func NewLocation(u *url.URL) (*ServiceLocation) {
	sl := ServiceLocation{*u}
	return &sl
}

func (s *Service) GetLocationsForVersion(v Version) ([]*ServiceLocation) {
	for i, value := range s.Versions {
		if (value.v.Matches(&v)) {
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
		if (value.v.Matches(&v)) {
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
