package ServiceRegistry

import (
	//	"fmt"
)

func (sr *ServiceRegistry) GenerateTestData() {

	a := sr.NewService("TestService")
	b := sr.NewService("SomeOtherService")

	//	Create the data structures

	a.AddServiceInstance(NewVersion("1.0.0"), NewLocation("http://10.1.0.1:1237/afvdafvdafv"))
	a.AddServiceInstance(NewVersion("1.0.0"), NewLocation("http://10.1.0.2:1238/qwefrwq"))
	a.AddServiceInstance(NewVersion("1.0.0"), NewLocation("http://10.1.0.3:1239/abdfbadfbadfbadfba"))
	a.AddServiceInstance(NewVersion("1.24.37"), NewLocation("http://localhost:123/qwerty"))
	b.AddServiceInstance(NewVersion("1.24.37"), NewLocation("https://kevin.valinor.local/blahz"))
	a.AddServiceInstance(NewVersion("2.24.37"), NewLocation("https://kevin.valinor.local/blahqwertyz"))
	a.AddServiceInstance(NewVersion("2.24.37"), NewLocation("https://kevin.valinor.local/sfvslfvsf"))

	//	fmt.Printf(a.Name + ": %d versions\n", len(a.Versions))
	//	fmt.Printf(b.Name + ": %d versions\n", len(b.Versions))
	//
	//	// Interrogate the data structures
	//	fmt.Println(sr.GetServiceWithName("TestService", false).GetLocationsForVersion(NewVersion("1.24.37"))[0].Location)
	//
	//	// Interrogate the data structures
	//	fmt.Println(sr.GetServiceWithName("TestService", false).GetLocationForVersion(Version("1.24.37")))


}
