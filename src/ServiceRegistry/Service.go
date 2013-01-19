package ServiceRegistry

type Service struct {
	name     string
	Versions []ServiceVersion
}

func (S *Service) AddServiceVersion(sv ServiceVersion) {
	S.Versions = append(S.Versions, sv);
}
