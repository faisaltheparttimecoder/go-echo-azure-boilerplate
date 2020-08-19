package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Check if the value is empty
func IsSettingEmpty(ev string) string {
	s := os.Getenv(ev)
	// Fail if there is no value for that environment variable
	if !DoesValueExists(s) {
		Fatalf("Mandatory parameter \"%s\" is missing from the environment variable or its empty", ev)
	}
	return s
}

// Trim spaces and provide the value
func DoesValueExists(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	return true
}

// Print the error in the stdout & and also return error
func PrintErrorAndReturn(s string, e error) error {
	err := fmt.Sprintf(s, e)
	Errorf(err)
	return e
}

// Convert the string and to number
func ConvertStringToNumber(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		Fatalf("Failed to convert the string \"%s\" to integer, err: %v", s, err)
	}
	return i
}