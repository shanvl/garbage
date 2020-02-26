// Package valid provides tools to validate some input and return a convenient error
package valid

// Check receives validation functions and returns an error if one of them returned "false"
func Check(validateFunctions ...func() (isValid bool, errorKey string, errorDescription string)) error {
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
