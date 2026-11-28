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

// Autohellas (Hertz GR) allow-lists the static address production and the proxy egress through.
// Production calls its production service directly; other cloud environments reach UAT through the
// proxy, which forwards /hertzgr/* to https://uatota.hertz.gr; local runs call UAT directly, from an
// allow-listed office address.
HertzGRBaseURL: [
	if #Meta.Environment.Type == "production" {"https://prdota01.hertz.gr/OTAServices/AutohellasService.svc"},
	if #Meta.Environment.Cloud == "local" {"https://uatota.hertz.gr/OTAServices/AutohellasService.svc"},
	"https://ai-booking-egress-proxy-jurknw5msa-zf.a.run.app/hertzgr/OTAServices/AutohellasService.svc",
][0]

// The rate codes Autohellas assigned us for Greece, one plan each: the code's 4th letter selects the
// brand (H Hertz, T Thrifty, D Dollar, F Firefly). Append new codes at the end: a code's position is
// its plan id, which price offers are renewed by. The ERP day charges are placeholders until decided.
//
// Production has none until Hertz GR goes live there: with no rate codes the broker returns no cars
// and calls Autohellas for nothing, so it can't be searched, offered or booked.
HertzGRRates: [
	if #Meta.Environment.Type == "production" {[]},
	[
		{Code: "ITRHABT", PlanName: "Standard", ErpDayCharge: 3.0},
		{Code: "ITRTABT", PlanName: "Standard", ErpDayCharge: 3.0},
		{Code: "ITRDABT", PlanName: "Standard", ErpDayCharge: 3.0},
	],
][0]

// HertzGRLogXML logs every Autohellas request and response, with the password and validation code
// masked. Autohellas certifies us from these exchanges before giving production access. Off now
// in every environment; turn it on again if Autohellas asks for exchanges.
HertzGRLogXML: false
