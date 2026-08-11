package billing

MonthlyReport: {
	Headers: {
		OfficeName:           "שם- משרד"
		AgentName:            "שם סוכן"
		DriverName:           "שם הנהג"
		ReservationCreatedAt: "תאריך ביצוע הזמנה"
		ReservationID:        "מספר הזמנה"
		VoucherDate:          "תאריך כירטוס"
		VoucherNumber:        "מספר שובר"
		AgentVoucherNumber:   "מספר שובר סוכן"
		PickupDate:           "תאריך איסוף"
		DropoffDate:          "תאריך החזרה"
		CountryCode:          "קוד מדינה"
		RentalDays:           "ימי השכרה"
		Currency:             "מטבע עסקה"
		NetPrice:             "מחיר נטו"
		FullCoverage:         "כיסוי מלא"
		TotalNetPrice:        "סה\"כ לתשלום נטו"
	}
	Styles: {
		HeaderBackgroundColor:    "D9D9D9"
		RefundRowBackgroundColor: "FFD8B1"
		TotalRowBackgroundColor:  "FFFF00"
		BorderColor:              "000000"
	}
}

if #Meta.Environment.Type != "production" && #Meta.Environment.Name != "staging" {
	Icount: {
		AccountID: 4,
		AgentsPaypageID: 3,
		CustomersPaypageID: 4,
	}
}

if #Meta.Environment.Type == "production" || #Meta.Environment.Name == "staging" {
	Icount: {
		AccountID: 4,
		AgentsPaypageID: 14,
		CustomersPaypageID: 18,
	}
}

Invoice:{
	PurchaseItemDescription: "עלות השכרה",
	ProfitItemDescription: "רווח השכרה",
	ProfitAndErpItemDescription: "רווח השכרה כולל כיסוי",
	CancellationPenaltyItemDescription: "דמי ביטול",
	NoShowPenaltyItemDescription: "דמי אי-הגעה",
	DiscountItemDescription: "הנחה",
}