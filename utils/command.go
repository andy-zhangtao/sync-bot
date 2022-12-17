package utils

import "strings"

func GrabCommand(msg string) string {
	_msg := strings.Split(msg, " ")
	if len(_msg) == 1 {
		return msg
	}

	result := strings.TrimSpace(strings.Trim(msg, _msg[0]))
	if result == "" {
		return msg
	}

	return result
}
