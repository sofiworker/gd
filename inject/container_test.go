package inject

import "testing"

func TestNew(t *testing.T) {
	container := New()
	module := container.NewModule("root")
	err := module.Provide("test", nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	subModule := module.SubModule("sub")
	err = subModule.Provide("test", nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
