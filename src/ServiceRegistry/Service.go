package ServiceRegistry

import (
	"net/url"
	"fmt"
)

type Service struct {
	Name     string
	Versions []ServiceVersion
}

func NewService(name string) (*Service) {
	s := new(Service)
	s.Name = name
	return s
}

func (S *Service) AddServiceVersion(v Version) {
	//sv := ServiceVersion.New()
	//S.Versions = append(S.Versions, sv)
}

func NewVersion(major int, minor int, micro int) (Version) {
	v := Version{major, minor, micro}
	return v
}

func NewLocation(u *url.URL) (ServiceLocation) {
	sl := ServiceLocation{*u}
	return sl
}

func (s *Service) getVersion(v Version) (*ServiceVersion) {
	for i, value := range s.Versions {
		if (value.v == v) {
			return &s.Versions[i]
		}
	}
	// Need to make a new ServiceVersion
	sv := ServiceVersion{v, nil}
	s.Versions = append(s.Versions, sv)
	return &s.Versions[0]
}

func (s *Service) AddServiceInstance(v Version, sl ServiceLocation) {
	sv := s.getVersion(v)
	sv.locations = append(sv.locations, sl)
	fmt.Printf("There are %d locations for version %d.%d.%d\n", len(sv.locations), v.Major, v.Minor, v.Micro)
	for _, value := range sv.locations {
		fmt.Printf(" - %v\n", value.location)
	}
}
