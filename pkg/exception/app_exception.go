package exception

import (
	"github.com/exgamer/gosdk/pkg/validation"
	"github.com/go-errors/errors"
	"github.com/gookit/validate"
	"net/http"
)

// AppException Модель данных для описания ошибки
type AppException struct {
	Code        int
	Error       error
	Context     map[string]any
	ServiceCode int
}

func NewAppException(code int, err error, context map[string]any) *AppException {
	return &AppException{code, err, context, 0}
}

func NewInternalServerAppException(err error, context map[string]any) *AppException {
	return &AppException{http.StatusInternalServerError, err, context, 0}
}

func NewValidationAppException(context map[string]any) *AppException {
	return &AppException{http.StatusUnprocessableEntity, errors.New("VALIDATION ERROR"), context, 0}
}

func NewValidationAppExceptionFromValidationErrors(validationErrors validate.Errors) *AppException {
	return NewValidationAppException(validation.ValidationErrorsAsMap(validationErrors))
}
