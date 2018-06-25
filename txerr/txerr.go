package txerr

type ErrorBuild interface {
	error
	Err() Error
}

type Error interface {
	error

	Op() string
	Kind() error
	Message() string
	Cause() error
}

// selective accessors
type (
	withOp       interface{ Op() string }
	withKind     interface{ Kind() error }
	withMessage  interface{ Message() string }
	withChild    interface{ Cause() error }
	withChildren interface{ Causes() []error }
	builder      interface{ Err() Error }
)

// FindErrWith returns the first error in the error tree, that matches the
// given predicate.
func FindErrWith(in error, pred func(err error) bool) error {
	var found error
	Iter(in, func(err error) bool {
		matches := pred(err)
		if matches {
			found = err
			return false
		}
		return true
	})

	return found
}

// Iter iterates the complete error tree call fn on each found error value.
// The user function fn can stop the iteration by returning false.
func Iter(in error, fn func(err error) bool) {
	doIter(in, fn)
}

func doIter(in error, fn func(err error) bool) bool {
	for {
		if in == nil {
			return true
		}

		if cont := fn(in); !cont {
			return cont
		}

		switch err := in.(type) {
		case builder:
			in = err.Err() // resolve ErrorBuild to Error and retry

		case withChild:
			in = err.Cause()

		case withChildren:
			for _, sub := range err.Causes() {
				if cont := doIter(sub, fn); !cont {
					return cont
				}
			}
			return true

		default:
			return true
		}
	}
}

func directMsg(in error) string {
	if err, ok := in.(withMessage); ok {
		return err.Message()
	}
	return ""
}
