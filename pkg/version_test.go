package pkg

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	vers := version
	if GetVersion() != vers {
		t.Errorf("GetVersion(): expected %s, actual %s", vers, GetVersion())
	}
}
