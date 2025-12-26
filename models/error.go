package models

type ComiketBackendError struct {
	Err        error
	StatusCode int
}

func (comiketBackendError ComiketBackendError) Status() int {
	return comiketBackendError.StatusCode
}

func (comiketBackendError ComiketBackendError) Error() string {
	return comiketBackendError.Err.Error()
}
