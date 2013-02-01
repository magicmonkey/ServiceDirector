package ServiceRegistry

import (
	"net/url"
	//	"fmt"
)

func (sr *ServiceRegistry) GenerateTestData() {

	a := sr.NewService("TestService")
	b := sr.NewService("SomeOtherService")

	//	Create the data structures

	var u *url.URL

	u, _ = url.Parse("http://10.1.0.1:1237/afvdafvdafv")
	a.AddServiceInstance(NewVersion("1.0.0"), NewLocation(u))

	u, _ = url.Parse("http://10.1.0.2:1238/qwefrwq")
	a.AddServiceInstance(NewVersion("1.0.0"), NewLocation(u))

	u, _ = url.Parse("http://10.1.0.3:1239/abdfbadfbadfbadfba")
	a.AddServiceInstance(NewVersion("1.0.0"), NewLocation(u))

	u, _ = url.Parse("http://localhost:123/qwerty")
	a.AddServiceInstance(NewVersion("1.24.37"), NewLocation(u))

	u, _ = url.Parse("https://kevin.valinor.local/blahz")
	b.AddServiceInstance(NewVersion("1.24.37"), NewLocation(u))

	u, _ = url.Parse("https://kevin.valinor.local/blahqwertyz")
	a.AddServiceInstance(NewVersion("2.24.37"), NewLocation(u))

	u, _ = url.Parse("https://kevin.valinor.local/sfvslfvsf")
	a.AddServiceInstance(NewVersion("2.24.37"), NewLocation(u))

	//	fmt.Printf(a.Name + ": %d versions\n", len(a.Versions))
	//	fmt.Printf(b.Name + ": %d versions\n", len(b.Versions))
	//
	//	// Interrogate the data structures
	//	fmt.Println(sr.GetServiceWithName("TestService", false).GetLocationsForVersion(NewVersion("1.24.37"))[0].Location)
	//
	//	// Interrogate the data structures
	//	fmt.Println(sr.GetServiceWithName("TestService", false).GetLocationForVersion(Version("1.24.37")))


}
