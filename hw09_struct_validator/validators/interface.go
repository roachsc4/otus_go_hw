package validators

type Validator interface {
	Init(validatorValue string) error
	Validate(valueToValidate interface{}) error
}
