package validators

import (
	"fmt"
	"strconv"
	"strings"
)

type MinValidatorError struct {
	Value    int
	MinValue int
}

func (e *MinValidatorError) Error() string {
	return fmt.Sprintf(
		"int %d must be greater or equal than %d",
		e.Value,
		e.MinValue)
}

type MaxValidatorError struct {
	Value    int
	MaxValue int
}

func (e *MaxValidatorError) Error() string {
	return fmt.Sprintf(
		"int %d must be equal to or lesser than %d",
		e.Value,
		e.MaxValue)
}

type IntInError struct {
	Value         int
	AllowedValues map[int]struct{}
}

func (e *IntInError) Error() string {
	return fmt.Sprintf(
		"int %d doesn't fit allowed set %v",
		e.Value,
		e.AllowedValues)
}

type MinValidator struct {
	minValue int
}

func (v *MinValidator) Init(validatorValue string) error {
	minValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	v.minValue = minValue

	return nil
}

func (v MinValidator) Validate(valueToValidate interface{}) error {
	intToValidate, ok := valueToValidate.(int)
	if !ok {
		return fmt.Errorf("unexpected Value %v", valueToValidate)
	}

	if intToValidate < v.minValue {
		return &MinValidatorError{
			Value:    intToValidate,
			MinValue: v.minValue,
		}
	}
	return nil
}

type MaxValidator struct {
	maxValue int
}

func (v *MaxValidator) Init(validatorValue string) error {
	minValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	v.maxValue = minValue

	return nil
}

func (v MaxValidator) Validate(valueToValidate interface{}) error {
	intToValidate, ok := valueToValidate.(int)
	if !ok {
		return fmt.Errorf("unexpected Value %v", valueToValidate)
	}

	if intToValidate > v.maxValue {
		return &MaxValidatorError{
			Value:    intToValidate,
			MaxValue: v.maxValue,
		}
	}
	return nil
}

type IntInValidator struct {
	allowedValues map[int]struct{}
}

func (v *IntInValidator) Init(validatorValue string) error {
	allowedValuesList := strings.Split(validatorValue, ",")
	v.allowedValues = make(map[int]struct{}, len(allowedValuesList))
	for _, value := range allowedValuesList {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("unexpected error: %w", err)
		}
		v.allowedValues[intValue] = struct{}{}
	}
	return nil
}

func (v IntInValidator) Validate(valueToValidate interface{}) error {
	intToValidate, ok := valueToValidate.(int)
	if !ok {
		return fmt.Errorf("unexpected Value %v", valueToValidate)
	}

	_, valueIsAllowed := v.allowedValues[intToValidate]
	if !valueIsAllowed {
		return &IntInError{
			Value:         intToValidate,
			AllowedValues: v.allowedValues,
		}
	}
	return nil
}
