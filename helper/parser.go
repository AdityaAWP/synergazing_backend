package helper

import (
	"encoding/json"
	"strconv"
)

func ParseStringSlice(jsonString string) ([]string, error) {
	if jsonString == "" {
		return []string{}, nil
	}

	var stringSlice []string
	err := json.Unmarshal([]byte(jsonString), &stringSlice)
	if err != nil {
		return nil, err
	}

	return stringSlice, nil
}

func StringToInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func StringToFloat(s string) (float64, error) {
	if s == "" {
		return 0.0, nil
	}
	return strconv.ParseFloat(s, 64)
}
