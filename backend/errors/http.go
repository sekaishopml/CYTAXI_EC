package errors

import "net/http"

func HTTPStatus(kind Kind) int {
	switch kind {
	case KindValidation:
		return http.StatusBadRequest
	case KindNotFound:
		return http.StatusNotFound
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindConflict:
		return http.StatusConflict
	case KindTimeout:
		return http.StatusGatewayTimeout
	case KindExternal:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

func KindFromHTTP(status int) Kind {
	switch {
	case status >= 500:
		return KindInternal
	case status == 409:
		return KindConflict
	case status == 404:
		return KindNotFound
	case status == 403:
		return KindForbidden
	case status == 401:
		return KindUnauthorized
	case status >= 400:
		return KindValidation
	default:
		return KindInternal
	}
}

func EncodeHTTP(err error) (int, any) {
	var e *Error
	if ok := As(err, &e); ok {
		return HTTPStatus(e.Kind), map[string]any{
			"error":   e.Message,
			"kind":    e.Kind,
			"op":      e.Op,
		}
	}
	return http.StatusInternalServerError, map[string]any{
		"error": "internal server error",
	}
}

func As(err error, target any) bool {
	if e, ok := err.(*Error); ok {
		*target.(**Error) = e
		return true
	}
	return false
}
