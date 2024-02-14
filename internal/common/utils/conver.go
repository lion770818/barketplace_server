package utils

import (
	"encoding/json"
	"strconv"
	"strings"
)

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

func ConverStringToFloat64(str string) (float64, error) {
	feetFloat, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	return feetFloat, err
}

func ConverJsonToMap(jsonStr string) (map[string]string, error) {
	byteArray, err := json.Marshal(jsonStr)
	if err != nil {
		return nil, err
	}

	var dataMap = make(map[string]string)
	err = json.Unmarshal(byteArray, &dataMap)
	if err != nil {
		return nil, err
	}
	return dataMap, nil
}
