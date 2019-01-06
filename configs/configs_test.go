package configs

import (
	"testing"
)

func TestLoadFile(t *testing.T) {
	var c configs
	Values.File =  "configs_test.yaml"
	LoadFile(&c)
	h := []string{"something", "crazy"}
	for i := range h {
		if i > len(c.Hostnames) - 1 || h[i] != c.Hostnames[i] {
			t.Fatalf("%v != %v", c.Hostnames, h)
		}
	}
	p := 8080
	if c.Port != p {
		t.Fatalf("%v != %v", c.Port, p)
	}
	
}
