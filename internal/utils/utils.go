package utils

import (
	"bufio"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

// generateUniqueId generates a 15 characters long string from lowercase letters and digits.
func GenerateUniqueId() string {
	const characters = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 15
	charactersLength := len(characters)
	uniqueId := make([]byte, length)

	for i := range uniqueId {
		randomIndex := rand.Intn(charactersLength)
		uniqueId[i] = characters[randomIndex]
	}

	return string(uniqueId)
}

// ValidateId returns true if id is 15 characters long made from lowercase letters and digits
func ValidateId(id string) bool {
	matched, _ := regexp.MatchString("^[a-z0-9]{15}$", id)
	return matched
}

// ParseEnv reads environmentral variables from the specified file
// and loads them into runtime.
func ParseEnv(filename string) error {
	envFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		envVar := strings.SplitN(scanner.Text(), "=", 2)
		os.Setenv(purifyString(envVar[0]), purifyString(envVar[1]))
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

// purifyString is a helper function that trims leading and trailing whitespace and removes quotation marks from a given string.
func purifyString(s string) string {
	s = strings.Trim(s, " ")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "'", "")
	return s
}

// NormalizeName converts underscore-separated parts of a name to spaces and capitalizes each part, preserving hyphens within each part.
func NormalizeName(name string) string {
	parts := strings.Split(name, "_")
	for i, part := range parts {
		subParts := strings.Split(part, "-")
		for j, subPart := range subParts {
			if len(subPart) > 0 {
				subParts[j] = strings.ToUpper(subPart[:1]) + strings.ToLower(subPart[1:])
			}
		}
		parts[i] = strings.Join(subParts, "-")
	}
	return strings.Join(parts, " ")
}
