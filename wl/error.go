package wl

// Internal error extracts the internal error cause from within the error returned by this package
func InternalError(e error) error {
	if err, ok := e.(combinedError); ok {
		return err[1]
	}
	return nil
}

type combinedError [2]error

func (c combinedError) Error() string {
	return c[0].Error() + ": " + c[1].Error()
}
func (c combinedError) Unwrap() error {
	return c[0]
}
