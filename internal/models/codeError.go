package models

type NoSuchLink struct {
	Message string
}

func (e *NoSuchLink) Error() string {
	return e.Message
}
