package utils

import (
	"regexp"
)

var (
	NamePattern           = regexp.MustCompile("^[a-zA-Z\\d_]{4,20}$")
	CapitalLettersPattern = regexp.MustCompile("[A-Z]")
	SmallLettersPattern   = regexp.MustCompile("[a-z]")
	DigitalPattern        = regexp.MustCompile("[0-9]")
	SpecialSymbolPattern  = regexp.MustCompile("[\\(\\)\\`~!@#\\$%\\^\\&\\*\\-\\+=\\_\\|\\{\\}\\[\\]:;'<>,.\\?/]")
)

func CheckInstancePwd(password string) bool {
	maxLength := 16
	minLength := 12

	lenPwd := len(password)
	if lenPwd < minLength || lenPwd > maxLength {
		return false
	}

	score := 0
	if CapitalLettersPattern.MatchString(password) {
		score += 1
	}

	if SmallLettersPattern.MatchString(password) {
		score += 1
	}

	if DigitalPattern.MatchString(password) {
		score += 1
	}

	if SpecialSymbolPattern.MatchString(password) {
		score += 1
	}

	if score >= 3 {
		return true
	}

	return false
}

func CheckName(name string) bool {
	return NamePattern.MatchString(name)
}
