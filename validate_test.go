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
func TestCheck(t *testing.T) {
	datos := []struct {
		name      string
		got       []string
		pattern   string
		wantError bool
	}{
		{"Valid1", []string{"kevin sauc", "kevin"}, "max=10", false},
		{"Invalid1", []string{"99999", "77777"}, "len=5", false},
		{"Invalid2", []string{"/kevin", "&gato"}, "basic", true},
	}
	k := New()
	for _, d := range datos {
		t.Run(d.name, func(t *testing.T) {
			if err := k.Target(d.pattern, d.got...).Ok(); !compare(err, d.wantError) {
				t.Errorf(err.Error())
			}
		})
	}
}

func TestStandardizeSpaces(t *testing.T) {
	samples := []struct {
		Name string
		Got  string
		Want string
	}{
		{"Input1", "key1 key2 key3", "key1 key2 key3"},
		{"Input1", "key1    key2     key3", "key1 key2 key3"},
		{"Input1", "       key1 key2 key3", "key1 key2 key3"},
		{"Input1", "key1 key2 key3        ", "key1 key2 key3"},
		{"Input1", "key1 key2      key3", "key1 key2 key3"},
	}
	for _, s := range samples {
		t.Run(s.Name, func(t *testing.T) {
			result := StandardizeSpaces(s.Got)
			if result != s.Want {
				t.Errorf(" Entrada: '%s', Se esperaba: '%s',Se obtuvo '%s'", s.Got, s.Want, result)
			}
		})
	}
}
func TestGetKeyArgs(t *testing.T) {
	type Got struct{ Targetkey, Key string }
	type Want struct {
		Args      string
		ExistArgs bool
	}
	samples := []struct {
		ID   string
		got  Got
		want Want
	}{
		{"Input1", Got{"max=20", "max="}, Want{"20", true}},
		{"Input3", Got{"max=", "max="}, Want{"", false}},
		{"Input4", Got{"", ""}, Want{"", false}},
	}
	for _, s := range samples {
		t.Run(s.ID, func(t *testing.T) {
			result, ok := GetKeyArgs(s.got.Targetkey, s.got.Key)
			if result != s.want.Args || ok != s.want.ExistArgs {
				t.Errorf("Se esperaba {%s,%v}, Se obtuvo {%s,%v}", s.want.Args, s.want.ExistArgs, result, ok)
			}
		})
	}
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
			if err := valid.Num(d.got, ""); !compare(err, d.wantError) {
				t.Errorf("Se esperaba: hayError=%v", d.wantError)
			}
		})
	}
}
func TestLen(t *testing.T) {
	datos := []struct {
		name      string
		got       string
		wantError bool
	}{
		{"Valid1", "kevin", false},
		{"Invalid1", "kevin1", true},
		{"Invalid2", "kevi", true},
	}
	valid := validate{}
	for _, d := range datos {
		t.Run(d.name, func(t *testing.T) {
			if err := valid.Len(d.got, "5"); !compare(err, d.wantError) {
				t.Errorf("Se esperaba: hayError=%v", d.wantError)
			}
		})
	}
}
func TestMaxMinLen(t *testing.T) {
	datos := []struct {
		name      string
		got       string
		wantError bool
	}{
		{"Valid1", "kevin", false},
		{"Valid2", "323323", false},
		{"Invalid1", "2123333", true},
		{"Invalid2", "Hola como estas", true},
	}
	valid := validate{}
	for _, d := range datos {
		t.Run(d.name, func(t *testing.T) {
			if err := valid.MaxLen(d.got, "6"); !compare(err, d.wantError) {
				t.Errorf("Se esperaba: hayError=%v", d.wantError)
			}
		})
	}
	for _, d := range datos {
		t.Run(d.name+"_min", func(t *testing.T) {
			if err := valid.MinLen(d.got, "7"); !compare(err, !d.wantError) {
				t.Errorf("Se esperaba: hayError=%v", d.wantError)
			}
		})
	}
}
func TestNoSpaces(t *testing.T) {
	datos := []struct {
		name      string
		got       string
		wantError bool
	}{
		{"Valid1", "kevin", false},
		{"Invalid1", "2123 333", true},
		{"Invalid2", " Hola", true},
		{"Invalid3", "Hola ", true},
		{"Invalid4", "Ho la", true},
	}
	valid := validate{}
	for _, d := range datos {
		t.Run(d.name, func(t *testing.T) {
			if err := valid.NoSpaces(d.got, ""); !compare(err, d.wantError) {
				t.Errorf("Se esperaba: hayError=%v", d.wantError)
			}
		})
	}
}
func TestValidWord(t *testing.T) {
	datos := []struct {
		name      string
		got       string
		wantError bool
	}{
		{"Valid1", "kevin", false},
		{"Invalid1", "2123 333", true},
		{"Invalid2", " Hola", true},
		{"Invalid3", "Hola ", true},
		{"Invalid4", "Ho la", true},
	}
	valid := validate{}
	for _, d := range datos {
		t.Run(d.name, func(t *testing.T) {
			if err := valid.NoSpaces(d.got, ""); !compare(err, d.wantError) {
				t.Errorf("Se esperaba: hayError=%v", d.wantError)
			}
		})
	}
}
func TestWord(t *testing.T) {
	datos := []struct {
		name      string
		got       string
		wantError bool
	}{
		{"Valid1", "hola como estas", false}, {"Valid2", "k", false}, {"Valid3", "", false},
		{"Invalid1", "hola como  estas", true}, {"Invalid2", " Hola", true}, {"Invalid3", "Hola ", true},
		{"Invalid4", "hola   como   estas", true},
	}
	valid := validate{}
	for _, d := range datos {
		t.Run(d.name, func(t *testing.T) {
			if err := valid.Words(d.got, ""); !compare(err, d.wantError) {
				t.Errorf("Se esperaba: hayError=%v", d.wantError)
			}
		})
	}
}

func TestBasic(t *testing.T) {
	datos := []struct {
		name      string
		got       string
		wantError bool
	}{
		{"Valid1", "Kevin Saucedo", false}, {"Valid3", "Jose    Días", false},
		{"Invalid1", "hol%", true}, {"Invalid2", "H*ola", true},
		{"Invalid4", "hola   como  :: estas", true},
	}
	valid := validate{}
	for _, d := range datos {
		t.Run(d.name, func(t *testing.T) {
			if err := valid.Basic(d.got, ""); !compare(err, d.wantError) {
				t.Errorf("Se esperaba: hayError=%v", d.wantError)
			}
		})
	}
}
