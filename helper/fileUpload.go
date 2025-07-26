package helper

import "strings"


func GetUrlFile(filepath string) string {
	if filepath == "" {
		return ""
	}
	return "/" + strings.ReplaceAll(filepath, "\\", "/")
}