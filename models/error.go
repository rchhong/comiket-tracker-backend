package models

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Err        error
	StatusCode int
}

func (statusError StatusError) Status() int {
	return statusError.StatusCode
}

func (statusError StatusError) Error() string {
	return statusError.Err.Error()
}

type ErrorResponse struct {
	Message string `json:"message"`
}
