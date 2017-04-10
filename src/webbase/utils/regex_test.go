package utils

import (
	"testing"

	"github.com/stretchr/testify.v2/require"
)

func TestCheckInstancePwd(t *testing.T) {
	pwd1 := "1234555"
	result1 := CheckInstancePwd(pwd1)
	require.Equal(t, false, result1)

	pwd2 := "ABCDEF"
	result2 := CheckInstancePwd(pwd2)
	require.Equal(t, false, result2)

	pwd3 := "abcdefeawfwe"
	result3 := CheckInstancePwd(pwd3)
	require.Equal(t, false, result3)

	pwd4 := "()`~!@#$%^&"
	result4 := CheckInstancePwd(pwd4)
	require.Equal(t, false, result4)

	pwd5 := "12345asdfc"
	result5 := CheckInstancePwd(pwd5)
	require.Equal(t, false, result5)

	pwd6 := "123456ADSSWF"
	result6 := CheckInstancePwd(pwd6)
	require.Equal(t, false, result6)

	pwd7 := "123456!@#$%^&"
	result7 := CheckInstancePwd(pwd7)
	require.Equal(t, false, result7)

	pwd8 := "asdfcADSSWF"
	result8 := CheckInstancePwd(pwd8)
	require.Equal(t, false, result8)

	pwd9 := "asdfgc!@#$%^&"
	result9 := CheckInstancePwd(pwd9)
	require.Equal(t, false, result9)

	pwd10 := "ADCSDW!@#$%^&"
	result10 := CheckInstancePwd(pwd10)
	require.Equal(t, false, result10)

	pwd11 := "12345ADWDddd"
	result11 := CheckInstancePwd(pwd11)
	require.Equal(t, true, result11)

	pwd12 := "123_-**AAW#dada22"
	result12 := CheckInstancePwd(pwd12)
	require.Equal(t, false, result12)

	pwd13 := "123_-**A#dada22"
	result13 := CheckInstancePwd(pwd13)
	require.Equal(t, true, result13)
}
