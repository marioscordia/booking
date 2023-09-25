package types

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// Define a new Validator type which contains a map of validation errors for our
// form fields.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
type Validator struct {
	NonFieldErrors []string
	FieldErrors map[string]string
} 
// Update the Valid() method to also check that the NonFieldErrors slice is
// empty.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}
// Create an AddNonFieldError() helper for adding error messages to the new
// NonFieldErrors slice.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
} 
// AddFieldError() adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
// Note: We need to initialize the map first, if it isn't already
// initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	} 
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
} 
// CheckField() adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
} 
// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
} 
// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
} 
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
func ConfirmPassword(value string, value2 string) bool{
	return value == value2
}
func CheckCategory(category []string) bool{
	return len(category) > 0
}
func CheckDates(start, end time.Time) bool{
	if start.After(end) || start.Before(time.Now()) || end.Before(time.Now()) {
		return false
	}
	return true
}