package kcheck

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	mumeros = "0123456789"
)

const (
	NUMERO    = "num"
	MAX_LEN   = "max="
	MIN_LEN   = "min="
	LEN       = "len="
	NO_SPACES = "no-spaces"
	WORDS     = "words"
	BASIC     = "basic"
)

const (
	InvalidKeyError = "Identificador invalido '%s'"
	NumError        = "Se esperaba un numero, obtuvimos %s, bad = '%s'; target ='%s'"
	MaxLenError     = "Len maximo: %d, obtuvimos %d; target= '%s'"
	MinLenError     = "Len Minimo: %d, obtuvimos %d; target= '%s'"
	LenError        = "Len: %d, obtuvimos %d; target= '%s'"
	NoSpacesError   = "No esta permitidos los espacios"
	WordsError      = "Solo palabras validas"
	BasicError      = "Solo palabras simples: carácter no permitido:'%s'"
)

type Validater interface {
	Ok() error
}
type ValidFunc func(tar string, args string) error
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
	v.Functions[MAX_LEN] = v.MaxLen
	v.Functions[MIN_LEN] = v.MinLen
	v.Functions[LEN] = v.Len
	v.Functions[NO_SPACES] = v.NoSpaces
	v.Functions[WORDS] = v.Words
	v.Functions[BASIC] = v.Basic
}

func StandardizeSpaces(s string) string {
	if len(s) == 0 {
		return s
	}
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.Trim(s, " ")
}

func (v *validate) MapPatterns(pattern string) {
	pattern = StandardizeSpaces(pattern)
	keys := strings.Split(pattern, " ")
	v.patternKeys = append(v.patternKeys, keys...)
}
func GetKeyArgs(targetKey, key string) (string, bool) {
	keyLen := len(key)
	targetKeyLen := len(targetKey)
	if keyLen == targetKeyLen {
		return "", false
	}
	if targetKeyLen < keyLen {
		return "", false
	}
	return targetKey[keyLen:], true
}
func ReadKeyIfIsValid(s string) (string, bool) {
	if strings.Contains(s, "=") {
		return strings.Split(s, "=")[0] + "=", true
	}
	return "", false
}

func (v *validate) ValidateTarges() error {
	for _, key := range v.patternKeys {
		f, ok := v.Functions[key]
		if ok {
			if err := v.ExecuteFunction(f, ""); err != nil {
				return err
			}
		} else {
			mykey, isValid := ReadKeyIfIsValid(key)
			if isValid {
				f, ok := v.Functions[mykey]
				if ok {
					args, ok := GetKeyArgs(key, mykey)
					if ok {
						if err := v.ExecuteFunction(f, args); err != nil {
							return err
						}
					} else {
						return fmt.Errorf(InvalidKeyError, key)
					}
				} else {
					return fmt.Errorf(InvalidKeyError, mykey)
				}
			} else {
				return fmt.Errorf(InvalidKeyError, key)
			}
		}
	}
	return nil
}

func (v *validate) ExecuteFunction(f ValidFunc, args string) error {
	for _, target := range v.Targets {
		if err := f(target, args); err != nil {
			return err
		}
	}
	return nil
}

func (v *validate) Num(tar string, _ string) error {
	l := len(tar)
	index := 0
	bad := ""
	for i, s := range tar {
		if i == 0 && (s == '-' || s == '+') {
			s = '0'
		}
		if !strings.Contains(mumeros, string(s)) {
			bad = string(s)
			break
		}
		index = i
	}
	if index == l-1 {
		return nil
	}
	return fmt.Errorf(NumError, tar, bad, tar)
}
func (v *validate) calclens(tar, arg string) (tarlen, lenArg int, err error) {
	if err = v.Num(arg, ""); err != nil {
		err = fmt.Errorf(InvalidKeyError+" "+err.Error(), MAX_LEN)
		return
	}
	lenArg, _ = strconv.Atoi(arg)
	tarlen = len(tar)
	err = nil
	return
}
func (v *validate) MaxLen(tar string, arg string) error {
	tarlen, lenArg, err := v.calclens(tar, arg)
	if err != nil {
		return err
	}
	if tarlen > lenArg {
		return fmt.Errorf(MaxLenError, lenArg, tarlen, tar)
	}
	return nil
}
func (v *validate) MinLen(tar string, arg string) error {
	tarlen, lenArg, err := v.calclens(tar, arg)
	if err != nil {
		return err
	}
	if tarlen < lenArg {
		return fmt.Errorf(MinLenError, tarlen, tarlen, tar)
	}
	return nil
}
func (v *validate) Len(tar string, arg string) error {
	tarlen, lenArg, err := v.calclens(tar, arg)
	if err != nil {
		return err
	}
	if tarlen != lenArg {
		return fmt.Errorf(LenError, lenArg, tarlen, tar)
	}
	return nil
}
func (v *validate) NoSpaces(tar string, _ string) error {
	if strings.Contains(tar, " ") {
		return fmt.Errorf(NoSpacesError)
	}
	return nil
}

func (v *validate) Words(tar string, _ string) error {
	if len(tar) == 0 {
		return nil
	}
	if tar[0] == ' ' || tar[len(tar)-1] == ' ' {
		return fmt.Errorf(WordsError)
	}
	if strings.Contains(tar, "  ") {
		return fmt.Errorf(WordsError)
	}
	return nil
}
func (v *validate) Basic(tar string, _ string) error {
	characters := "!\"#$%&'()*+,./:;<=>?@[\\]^_{|}~"
	for _, t := range tar {
		if strings.ContainsRune(characters, t) {
			return fmt.Errorf(BasicError, string(t))
		}
	}
	return nil
}

func (v *validate) Ok() error {
	if len(v.patternKeys) == 0 || len(v.Targets) == 0 {
		return nil
	}
	return v.ValidateTarges()
}
