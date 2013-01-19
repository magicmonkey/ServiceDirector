/**
 * Created with IntelliJ IDEA.
 * User: kevin
 * Date: 19/01/13
 * Time: 14:42
 * To change this template use File | Settings | File Templates.
 */
package ServiceRegistry

type ServiceVersion struct {
	major    int
	minor    int
	micro    int
	location []ServiceLocation
}
