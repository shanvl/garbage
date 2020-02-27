// Package valid provides tools to validate some input and return a convenient error
package valid

// CheckErrors receives validation functions and returns a convenient error if one of them returned "false"
func CheckErrors(validationFuncs ...func() (isValid bool, errorKey string, errorDescription string)) error {
	err := &ErrValidation{}
	for _, f := range validationFuncs {
		if isValid, errorKey, errorDescription := f(); !isValid {
			err.add(errorKey, errorDescription)
		}
	}
	if len(err.fields) > 0 {
		return err
	}
	return nil
}
