package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// GenerateID generates a random hex ID
func GenerateID(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^[\+]?[\d\s\-\(\)]{10,15}$`)
	return phoneRegex.MatchString(phone)
}

// SanitizeString removes extra spaces and trims string
func SanitizeString(s string) string {
	// Remove extra spaces
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// IsEmpty checks if a value is empty
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
}

// FormatTime formats time to string with specific layout
func FormatTime(t time.Time, layout string) string {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return t.Format(layout)
}

// ParseTime parses string to time with specific layout
func ParseTime(timeStr, layout string) (time.Time, error) {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return time.Parse(layout, timeStr)
}

// ContainsString checks if slice contains a string
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsInt checks if slice contains an integer
func ContainsInt(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// StringToPointer converts string to string pointer
func StringToPointer(s string) *string {
	return &s
}

// IntToPointer converts int to int pointer
func IntToPointer(i int) *int {
	return &i
}

// BoolToPointer converts bool to bool pointer
func BoolToPointer(b bool) *bool {
	return &b
}

// PointerToString safely converts string pointer to string
func PointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// PointerToInt safely converts int pointer to int
func PointerToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// PointerToBool safely converts bool pointer to bool
func PointerToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// SliceString creates a string slice from varargs
func SliceString(items ...string) []string {
	return items
}

// SliceInt creates an int slice from varargs
func SliceInt(items ...int) []int {
	return items
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Clamp clamps a value between min and max
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ToSnakeCase converts CamelCase to snake_case
func ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// ToCamelCase converts snake_case to CamelCase
func ToCamelCase(str string) string {
	words := strings.Split(str, "_")
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			words[i] = strings.ToUpper(string(words[i][0])) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}

// Truncate truncates a string to a maximum length
func Truncate(str string, maxLength int) string {
	if len(str) <= maxLength {
		return str
	}
	return str[:maxLength]
}

// PadLeft pads a string to the left with a character
func PadLeft(str string, totalLength int, padChar rune) string {
	padLength := totalLength - len(str)
	if padLength <= 0 {
		return str
	}
	return strings.Repeat(string(padChar), padLength) + str
}

// PadRight pads a string to the right with a character
func PadRight(str string, totalLength int, padChar rune) string {
	padLength := totalLength - len(str)
	if padLength <= 0 {
		return str
	}
	return str + strings.Repeat(string(padChar), padLength)
}

// ReverseString reverses a string
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// SafeDivide performs division with zero check
func SafeDivide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}
