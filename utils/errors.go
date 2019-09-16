package utils

type FatalError struct {
	Message string
}

func (e FatalError) Error() string {
	return e.Message
}
