package kcheck

import (
	"fmt"
	"strings"
)

const (
	mumeros = "0123456789"
)

const (
	NUMERO = "num"
)

const (
	InvalidKeyError = "Identificador invalido '%s'"
	NumError        = "Se esperaba un numero, obtubimos %s, bad = '%s'"
)

type Validater interface {
	Ok() error
}
type ValidFunc func(tar string) error
type MapFunc map[string]ValidFunc

type validate struct {
	patternKeys []string
	Targets     []string
	Functions   MapFunc
}

func newValidate(pattern string, targets []string) Validater {
	v := &validate{Targets: targets}
	v.Functions = make(MapFunc)
	v.patternKeys = make([]string, 0)
	v.Register()
	v.MapPatterns(pattern)
	return v
}
func (v *validate) Register() {
	v.Functions[NUMERO] = v.Num
}

func (v *validate) MapPatterns(pat string) {
	pat = strings.Trim(pat, " ")
	keys := strings.Split(pat, " ")
	v.patternKeys = append(v.patternKeys, keys...)
}

func (v *validate) ValidateTarges() error {
	for _, key := range v.patternKeys {
		if f, ok := v.Functions[key]; !ok {
			return fmt.Errorf(InvalidKeyError, key)
		} else {
			for _, tar := range v.Targets {
				if err := f(tar); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *validate) Num(tar string) error {
	l := len(tar)
	index := 0
	bad := ""
	for i, s := range tar {
		if i == 0 && (s == '-' || s == '+') {
			s = '0'
		}
		if strings.Index(mumeros, string(s)) == -1 {
			bad = string(s)
			break
		}
		index = i
	}
	if index == l-1 {
		return nil
	}
	return fmt.Errorf(NumError, tar, bad)
}

func (v *validate) Ok() error {
	if len(v.patternKeys) == 0 || len(v.Targets) == 0 {
		return nil
	}
	return v.ValidateTarges()
}
