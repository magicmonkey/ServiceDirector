package ServiceRegistry

import "net/url"

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

func (s *Service) AddServiceInstance(v Version, sl ServiceLocation) {

}
