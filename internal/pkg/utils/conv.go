package utils

import "strconv"

func StrToUint(s string) (uint, error) {
	var i uint64

	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(i), nil
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func UintToStr(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

func BoolToStr(i bool) string {
	return strconv.FormatBool(i)
}
