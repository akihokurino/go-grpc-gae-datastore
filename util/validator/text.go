package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ValidateTextRange(s string, min int, max int) error {
	length := utf8.RuneCountInString(s)
	if length < min || max < length {
		return fmt.Errorf("range error, got %s", s)
	}
	return nil
}

func ValidateEmail(s string) error {
	rep := regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	if !rep.MatchString(s) {
		return fmt.Errorf("email invalid error, got %s", s)
	}
	return nil
}

func ValidatePhoneNumber(s string) error {
	rep1 := regexp.MustCompile(`^0\d\d{4}\d{4}$`)
	rep2 := regexp.MustCompile(`^(050|070|080|090)\d{4}\d{4}$`)
	rep3 := regexp.MustCompile(`^0120\d{3}\d{3}$`)
	if !rep1.MatchString(s) &&
		!rep2.MatchString(s) &&
		!rep3.MatchString(s) {
		return fmt.Errorf("phoneNumber invalid error, got %s", s)
	}
	return nil
}

func ValidatePostalCode(s string) error {
	rep := regexp.MustCompile(`^\d{7}$`)
	if !rep.MatchString(s) {
		return fmt.Errorf("postalCode invalid error, got %s", s)
	}
	return nil
}

func ValidateHiragana(s string) error {
	isHiragana := true
	kt := strings.Replace(s, "　", "", -1)
	kt = strings.Replace(kt, " ", "", -1)
	kt = strings.Replace(kt, "ー", "", -1)
	kt = strings.Replace(kt, "-", "", -1)

	for _, r := range kt {
		if !unicode.In(r, unicode.Hiragana) {
			isHiragana = false
			break
		}
	}

	if !isHiragana {
		return fmt.Errorf("not hiragana error, got %s", s)
	}

	return nil
}
