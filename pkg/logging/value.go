package logging

func ErrorValue(err error) any {
	return ErrorMessage(err)
}
