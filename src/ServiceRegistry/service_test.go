package ServiceRegistry

import "testing"

func TestPass(t *testing.T) {
	sr := ServiceRegistry{}
	sr.GenerateTestData()
}

func TestCreateNewService(t *testing.T) {
	sr := ServiceRegistry{}
	if sr.GetServiceWithName("Test", true) == nil {
		t.Fail()
	}
}

func TestDoNotCreateNewService(t *testing.T) {
	sr := ServiceRegistry{}
	if sr.GetServiceWithName("Test", false) != nil {
		t.Fail()
	}
}
