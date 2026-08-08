package model

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotFound   Error = "Resource not found"
	ErrInternal   Error = "Internal server errror"
	ErrBadRequest Error = "Bad request"
	ErrForbidden  Error = "Forbidden"
)
