package validators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type StringLengthError struct {
	Value          string
	RequiredLength int
}

func (e *StringLengthError) Error() string {
	return fmt.Sprintf(
		"string length is %d, but length %d is required",
		len(e.Value),
		e.RequiredLength)
}

type RegexpError struct {
	Value               string
	RegexpPatternString string
}

func (e *RegexpError) Error() string {
	return fmt.Sprintf(
		"string %s doesn't fit pattern %s",
		e.Value,
		e.RegexpPatternString)
}

type StringInError struct {
	Value         string
	AllowedValues map[string]struct{}
}

func (e *StringInError) Error() string {
	return fmt.Sprintf(
		"string %s doesn't fit allowed set %v",
		e.Value,
		e.AllowedValues)
}

type StringLengthValidator struct {
	requiredLength int
}

func (v *StringLengthValidator) Init(validatorValue string) error {
	requiredLength, err := strconv.Atoi(validatorValue)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	v.requiredLength = requiredLength

	return nil
}

func (v StringLengthValidator) Validate(valueToValidate interface{}) error {
	stringToValidate := fmt.Sprintf("%v", valueToValidate)

	if len(stringToValidate) != v.requiredLength {
		return &StringLengthError{
			Value:          stringToValidate,
			RequiredLength: v.requiredLength,
		}
	}
	return nil
}

type RegexpValidator struct {
	regexpPattern *regexp.Regexp
}

func (v *RegexpValidator) Init(validatorValue string) error {
	v.regexpPattern = regexp.MustCompile(validatorValue)
	return nil
}

func (v RegexpValidator) Validate(valueToValidate interface{}) error {
	stringToValidate := fmt.Sprintf("%v", valueToValidate)
	matched := v.regexpPattern.MatchString(stringToValidate)
	if !matched {
		return &RegexpError{
			Value:               stringToValidate,
			RegexpPatternString: v.regexpPattern.String(),
		}
	}
	return nil
}

type StringInValidator struct {
	allowedValues map[string]struct{}
}

func (v *StringInValidator) Init(validatorValue string) error {
	allowedValuesList := strings.Split(validatorValue, ",")
	v.allowedValues = make(map[string]struct{}, len(allowedValuesList))
	for _, value := range allowedValuesList {
		v.allowedValues[value] = struct{}{}
	}
	return nil
}

func (v StringInValidator) Validate(valueToValidate interface{}) error {
	stringToValidate := fmt.Sprintf("%v", valueToValidate)
	_, valueIsAllowed := v.allowedValues[stringToValidate]
	if !valueIsAllowed {
		return &StringInError{
			Value:         stringToValidate,
			AllowedValues: v.allowedValues,
		}
	}
	return nil
}
