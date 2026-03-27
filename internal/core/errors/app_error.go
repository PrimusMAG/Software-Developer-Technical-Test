package errors

type AppError struct {
	Message string
	Code    string
}

func (e AppError) Error() string {
	return e.Message
}
