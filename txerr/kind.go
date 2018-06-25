package txerr

// Is checks if any error in the error tree matches `kind`.
func Is(kind error, in error) bool {
	return FindKind(in, kind) != nil
}

func directKind(in error) error {
	if err, ok := in.(withKind); ok {
		return err.Kind()
	}
	return nil
}

// GetKind returns the first error kind found in the error tree.
func GetKind(in error) error {
	err := FindErrWith(in, func(err error) bool {
		if err, ok := in.(withKind); ok {
			return err.Kind() != nil
		}
		return false
	})

	if err == nil {
		return nil
	}
	return err.(withKind).Kind()
}

// FindKind returns the first error that matched `kind`.
func FindKind(in error, kind error) error {
	return FindErrWith(in, func(in error) bool {
		if err, ok := in.(withKind); ok {
			return err.Kind() == kind
		}
		return false
	})
}
