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

func NewAppException(code int, err error, context map[string]any, serviceCode int) *AppException {
	return &AppException{code, err, context, serviceCode}
}

func NewValidationAppException(context map[string]any, serviceCode int) *AppException {
	return &AppException{http.StatusUnprocessableEntity, errors.New("VALIDATION ERROR"), context, serviceCode}
}

func NewValidationAppExceptionFromValidationErrors(validationErrors validate.Errors, serviceCode int) *AppException {
	return NewValidationAppException(validation.ValidationErrorsAsMap(validationErrors), serviceCode)
}
