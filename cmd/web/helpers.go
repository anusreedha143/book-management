package main

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Define a Validator struct to hold validation errors
type Validator struct {
	Errors map[string]string
}

// NewValidator creates a new Validator instance
func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if there are no errors in the map
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map (if the key doesn't already exist)
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message if a validation condition is NOT met
// This is the "magic" method that makes your handler code clean
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// --- GENERIC VALIDATION HELPERS ---

// NotBlank returns true if a value is not empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if a value contains no more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinInt returns true if a value is greater than n
func MinInt(value int, n int) bool {
	return value > n
}

// Matches returns true if a value matches a provided regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
