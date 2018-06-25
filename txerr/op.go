package txerr

// IsOp checks if any error in the error tree is caused by `op`.
func IsOp(in error, op string) bool {
	return FindOp(in, op) != nil
}

func directOp(in error) string {
	if err, ok := in.(withOp); ok {
		return err.Op()
	}
	return ""
}

// GetOp returns the first errors it's Op value.
func GetOp(in error) string {
	err := FindErrWith(in, func(err error) bool {
		if err, ok := in.(withOp); ok {
			return err.Op() != ""
		}
		return false
	})

	if err == nil {
		return ""
	}
	return err.(withOp).Op()
}

// FindIp returns the first error with the given `op` value.
func FindOp(in error, op string) error {
	return FindErrWith(in, func(in error) bool {
		if err, ok := in.(withOp); ok {
			return err.Op() == op
		}
		return false
	})
}
