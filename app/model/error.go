package model

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotFound Error = "Resource not found"
	ErrInternal       = "Internal server errror"
	ErrBadRequest     = "Bad request"
)
