package booking

HertzErpDayChargeUS: 3.0
HertzErpDayChargeCA: 7.0
FlexErpDayCharge: 3.0
AvanceErpDayChargeStandard: 9.0
AvanceErpDayChargeZeroExcess: 3.0
SellBrokerErpAsOwn: true
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

// Autohellas (Hertz GR): production calls its production service, every other environment calls
// UAT directly.
HertzGRBaseURL: [
	if #Meta.Environment.Type == "production" {"https://prdota01.hertz.gr/OTAServices/AutohellasService.svc"},
	"https://uatota.hertz.gr/OTAServices/AutohellasService.svc",
][0]

// The rate codes Autohellas assigned us for Greece, one plan each: the code's 4th letter selects the
// brand (H Hertz, T Thrifty, D Dollar, F Firefly). Append new codes at the end: a code's position is
// its plan id, which price offers are renewed by. The ERP day charges are placeholders until decided.
HertzGRRates: [
	{Code: "ITRHABT", PlanName: "Standard", ErpDayCharge: 3.0},
	{Code: "ITRTABT", PlanName: "Standard", ErpDayCharge: 3.0},
	{Code: "ITRDABT", PlanName: "Standard", ErpDayCharge: 3.0},
]

// HertzGRLogXML logs every Autohellas request and response, with the password and validation code
// masked. Autohellas certifies us from these exchanges before giving production access. Off now
// in every environment; turn it on again if Autohellas asks for exchanges.
HertzGRLogXML: false
