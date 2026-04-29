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

type MonthlyReportData struct {
	ContactName string
}

var MonthlyReportTemplate = Template[MonthlyReportData]{
	name: "monthly_report",
}
