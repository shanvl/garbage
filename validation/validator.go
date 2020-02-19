// package validation provides tools to validate some input and return a convenient error
package validation

// Validator is a general validation interface
type Validator interface {
	// Validate receives validation functions and returns an error if one of them returned "false"
	Validate(validateFunctions ...func() (isValid bool, errorKey string, errorDescription string)) error
}

type validator struct{}

// Validate receives validation functions and returns an error if one of them returned "false"
func (v *validator) Validate(validateFunctions ...func() (isValid bool, errorKey string, errorDescription string)) error {
	err := &ErrValidation{}
	for _, f := range validateFunctions {
		if isValid, errorKey, errorDescription := f(); !isValid {
			err.add(errorKey, errorDescription)
		}
	}
	if len(err.fields) > 0 {
		return err
	}
	return nil
}

// NewValidator returns an instance of Validator
func NewValidator() Validator {
	return &validator{}
}
