package model

type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return "Not Found"
}
