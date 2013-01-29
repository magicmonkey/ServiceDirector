package ServiceRegistry

import "testing"

func TestPass(t *testing.T) {
}

func TestFail(t *testing.T) {
	t.Failf("Example fail")
}

func TestError(t *testing.T) {
	t.Errorf("Example error")
}
