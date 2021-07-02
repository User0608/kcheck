package kcheck

import (
	"testing"
)

func compare(err error, wantError bool) bool {
	if err == nil && !wantError {
		return true
	}
	if err != nil && wantError {
		return true
	}
	return false
}

func TestNum(t *testing.T) {
	datos := []struct {
		name      string
		got       string
		wantError bool
	}{
		{"ValidNum1", "3820", false},
		{"ValidNum2", "-3820", false},
		{"ValidNum3", "+3820", false},
		{"InvalidNum1", "hola", true},
		{"InvalidNum2", "- 332", true},
		{"InvalidNum3", "+ 443", true},
		{"InvalidNum4", "343+", true},
		{"InvalidNum4", " 3820 ", true},
	}
	valid := validate{}
	for _, d := range datos {
		t.Run(d.name, func(t *testing.T) {
			if err := valid.Num(d.got); !compare(err, d.wantError) {
				t.Errorf("Se esperaba: hayError=%v", d.wantError)
			}
		})
	}
}
