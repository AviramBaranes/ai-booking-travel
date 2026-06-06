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

type VoucherEmailData struct {
	VoucherNumber string
}

var VoucherEmailTemplate = Template[VoucherEmailData]{
	name: "send_voucher",
}

type CancellationEmailData struct {
	BookingReferenceID string
	DriverFullName     string
}

var CancellationEmailTemplate = Template[CancellationEmailData]{
	name: "cancellation_email",
}

type LateCancellationAlertEmailData struct {
	ReservationID       int64
	BrokerReservationID string
	AgentLabel          string
	OfficeLabel         *string
	OrganizationLabel   *string
}

var LateCancellationAlertEmailTemplate = Template[LateCancellationAlertEmailData]{
	name: "late_cancellation_alert",
}

type NewOrderEmailData struct {
	BookingReferenceID string
	DriverFullName     string
}

var NewOrderEmailTemplate = Template[NewOrderEmailData]{
	name: "new_order_email",
}

type OpenOrderAlertEmailData struct {
	BookingReferenceID string
	DriverFullName     string
}

var OpenOrderAlertEmailTemplate = Template[OpenOrderAlertEmailData]{
	name: "open_order_alert_email",
}
