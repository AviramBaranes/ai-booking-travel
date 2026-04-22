package email

type Template[T any] struct {
	name string
}

type CriticalErrorData struct {
	Message string
}

var CriticalErrorTemplate = Template[CriticalErrorData]{
	name: "critical_error",
}
