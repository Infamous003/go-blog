package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator contains a map of strings to strings representing errors
type Validator struct {
	Errors map[string]string
}

// Returns an initialized Validator
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Returns a boolean value. False is invalid, meaning there are errors, and true is valid
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// Adds an error to Validator with the provided key and value, if it doesn't already exists
func (v *Validator) AddError(key string, value string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = value
	}
}

// Evaluates the expression, and if !ok, adds an erroor to Validator
func (v *Validator) Check(ok bool, key string, value string) {
	if !ok {
		v.AddError(key, value)
	}
}

// Unique returns whether the values in a slice are unique or not
// Being a generic, it works with comparables like int, strings, floats, etc
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(uniqueValues) == len(values)
}

// Matches compares `value` against rx
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
