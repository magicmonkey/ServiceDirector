package ServiceRegistry

import "testing"

func TestPass(t *testing.T) {
	sr := ServiceRegistry{}
	sr.GenerateTestData()
}

func TestCreateNewService(t *testing.T) {
	sr := ServiceRegistry{}
	if _, created := sr.GetServiceWithName("Test", true); created {
		return
	}
	t.Fail()
}

func TestDoNotCreateNewService(t *testing.T) {
	sr := ServiceRegistry{}
	if _, created := sr.GetServiceWithName("Test", false); !created {
		return
	}
	t.Fail()
}
