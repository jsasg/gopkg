package respond

import (
	"reflect"

	"github.com/jsasg/gopkg/respond/errors"
)

type Option func(rep *responser)

func WithHttpStatus(status int) Option {
	return func(rep *responser) {
		rep.HttpStatus = status
	}
}

func WithCode(code int) Option {
	return func(rep *responser) {
		rep.Code = code
	}
}

func WithData(data interface{}) Option {
	return func(rep *responser) {
		if reflect.ValueOf(data).IsZero() {
			if reflect.TypeOf(data).Kind() == reflect.Array || reflect.TypeOf(data).Kind() == reflect.Slice {
				rep.Data = []struct{}{}
			}
			return
		}

		rep.Data = data
	}
}

func WithDataIgnoreEncrypt(ignore bool) Option {
	return func(rep *responser) {
		rep.ignoreEncrypt = ignore
	}
}

func WithMessage(msg string) Option {
	return func(rep *responser) {
		rep.Message = msg
	}
}

func WithError(err error) Option {
	return func(rep *responser) {
		rep.Error = err

		var rErr *errors.RespondError
		if errors.As(err, &rErr) {
			err = rErr
		}
		rep.Message = err.Error()
	}
}
