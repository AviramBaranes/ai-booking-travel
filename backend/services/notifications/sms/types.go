package sms

type smsRequest struct {
	SMS smsRequestSMS `json:"sms"`
}

type smsRequestSMS struct {
	User         smsRequestSMSUser          `json:"user"`
	Source       string                     `json:"source"`
	Destinations []smsRequestSMSDestination `json:"destinations"`
	Message      string                     `json:"message"`
}

type smsRequestSMSUser struct {
	Username string `json:"username"`
}

type smsRequestSMSDestination struct {
	Phone string `json:"phone"`
}
