package frozen

import "github.com/go-errors/errors"

// Wrap errors, returning a nil error if errors.Wrap returns a nil *Error.
func errorsWrap(e interface{}, skip int) error { // nolint:unparam
	// nolint:revive
	if err := errors.Wrap(e, skip+1); err != nil {
		return err
	}
	return nil
}

type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

const (
	WTF           = InternalError("should never be called!")
	Unimplemented = InternalError("not implemented")
)
