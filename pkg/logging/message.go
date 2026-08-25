package logging

func ErrorMessage(err error) string {
	return err.Error()
}
