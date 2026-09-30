package booking

HertzErpDayChargeUS: 3.0
HertzErpDayChargeCA: 7.0
FlexErpDayCharge: 3.0
AvanceErpDayChargeStandard: 9.0
AvanceErpDayChargeZeroExcess: 3.0
MarkUpGross: 35.0
MarkUpNet:   16.0

CID:       "aibookingtravel"
User:      "accounting"
AccountID: 4

// Wheelsys allow-lists the static NAT address production egresses through, so only production
// calls it directly. Every other environment goes via the Cloud Run proxy, which forwards
// /wheelsys/* to the same host and exits through that same address.
AvanceBaseURL: [
	if #Meta.Environment.Type == "production" {"https://endpoint.wheelsys.io"},
	"https://ai-booking-egress-proxy-jurknw5msa-zf.a.run.app/wheelsys",
][0]
