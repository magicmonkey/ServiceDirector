/**
 * Created with IntelliJ IDEA.
 * User: kevin
 * Date: 19/01/13
 * Time: 14:41
 * To change this template use File | Settings | File Templates.
 */
package ServiceRegistry

type Service struct {
	name     string
	Versions []ServiceVersion
}

func (S *Service) AddServiceVersion(sv ServiceVersion) {
	S.Versions = append(S.Versions, sv);
}
