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
	for _, value := range s.Versions {
		if (value.v == v) {
			fmt.Printf("Got a match for version %d.%d.%d\n", v.Major, v.Minor, v.Micro)
			return &value
		}
	}
	fmt.Printf("No match for version %d.%d.%d\n", v.Major, v.Minor, v.Micro)
	// Need to make a new ServiceVersion
	sv := ServiceVersion{v, nil}
	s.Versions = append(s.Versions, sv)
	return &sv
}

func (s *Service) AddServiceInstance(v Version, sl ServiceLocation) {
	sv := s.getVersion(v)
	sv.locations = append(sv.locations, sl)
	fmt.Printf("There are %d locations\n", len(sv.locations))
}
