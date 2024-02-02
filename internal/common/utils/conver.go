package utils

import "strconv"

// uint64 to string
func ConvertUintToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}

// string to uint64
func ConvertStringToUint(id string) (uint64, error) {
	return strconv.ParseUint(id, 10, 64)
}

// string to int
func ConvertStringToInt(id string) (int, error) {
	return strconv.Atoi(id)
}
