package validator

import (
	"fmt"
	"regexp"

	"github.com/slok/ragnarok/api"
)

var (
	validLabelChars = `^[a-zA-Z0-9\.\-_]+$`
	validLabelRegex = regexp.MustCompile(validLabelChars)
)

// errorIfEmptyStr will return an error if the string is empty
func errorIfEmptyStr(str string) error {
	if str == "" {
		return fmt.Errorf("the string is empty")
	}
	return nil
}

// errorIfBiggerStr will return an error if the string is bigger than the max length.
func errorIfBiggerStr(str string, max int) error {
	if len(str) > max {
		return fmt.Errorf("the '%s' string is bigger than %d characters", str, max)
	}
	return nil
}

// errorIfLesserStr will return an error if the string is smaller than the min length.
func errorIfLesserStr(str string, min int) error {
	if len(str) < min {
		return fmt.Errorf("the '%s' string is smaller than %d characters", str, min)
	}
	return nil
}

// errorIfInvalidLabel will check if the label is valid
func errorIfvalidLabelStr(str string) error {
	if !validLabelRegex.MatchString(str) {
		return fmt.Errorf("the '%s' is not a valid formatted string", str)
	}
	return nil
}

// errorIfInvalidLabel will check if the label is valid
func errorIfNoLabelKeys(labels map[string]string, keys []string) []error {
	errors := []error{}

	for _, key := range keys {

		if v, ok := labels[key]; !ok {
			errors = append(errors, fmt.Errorf("missing '%s' label key", key))
		} else {
			if v == "" {
				errors = append(errors, fmt.Errorf("'%s' label key value cant be empty", key))
			}
		}
	}
	return errors
}

// errorIfWrongTypeMeta will check if the object has a correct type meta
func errorIfWrongTypeMeta(obj api.Object) error {
	// TODO: check also a valid type.
	if obj.GetObjectKind() == "" || obj.GetObjectVersion() == "" {
		return fmt.Errorf("object type metadata is missing")
	}
	return nil
}
