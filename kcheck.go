package kcheck

type KChecker interface {
	Traget(string, ...string) Validater
}
type kcheck struct {
	val validate
}

func New() KChecker {
	return &kcheck{}
}
func (k *kcheck) Traget(pattern string, targets ...string) Validater {
	v := newValidate(pattern, targets)
	return v
}
