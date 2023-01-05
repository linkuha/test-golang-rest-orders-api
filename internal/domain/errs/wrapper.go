package errs

import "errors"

type CustomErrorWrapper struct {
	Message string
	Code    int
	Err     error
}

func (err CustomErrorWrapper) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

func (err CustomErrorWrapper) Unwrap() error {
	return err.Err
}

func (err CustomErrorWrapper) ErrorStack() (res string) {
	var ew CustomErrorWrapper
	if errors.As(err.Err, &ew) {
		res += err.Message + ": " + ew.ErrorStack()
	} else {
		res += err.Message + ": " + err.Err.Error() + ";"
	}
	return res
}

func (err CustomErrorWrapper) Dig() CustomErrorWrapper {
	var ew CustomErrorWrapper
	if errors.As(err.Err, &ew) {
		// Recursively digs until wrapper error is not CustomErrorWrapper
		return ew.Dig()
	}
	return err
}

func NewErrorWrapper(code int, err error, message string) error {
	return CustomErrorWrapper{
		Message: message,
		Code:    code,
		Err:     err,
	}
}
