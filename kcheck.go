package kcheck

type KChecker interface {
	Target(string, ...string) Validater
}

type kcheck struct {
	//val validate
}

func New() KChecker {
	return &kcheck{}
}

// Target pattern: patrones de validacion, targets, objetivos
func (k *kcheck) Target(pattern string, targets ...string) Validater {
	v := newValidate(pattern, targets)
	return v
}
