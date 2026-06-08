package notifications

AccountsEmailFrom: "accounting@aibookingtravel.com"
AccountsEmailFromName: "BookingTravel Accounts"
ReservationsEmailFrom: "reservations@aibookingtravel.com"
ReservationsEmailFromName: "AI Booking Travel"

SMSUsername: "sogomatic"
SMSSenderName: "Sogomatic"

GotenbergURL: [
	if #Meta.Environment.Type == "development" && #Meta.Environment.Cloud == "local" { "http://localhost:8080" },
	if #Meta.Environment.Name == "staging" { "https://gotenberg-217822056800.me-west1.run.app" },
	if #Meta.Environment.Type == "production" { "https://gotenberg-217822056800.me-west1.run.app" },
	"https://gotenberg-217822056800.me-west1.run.app",
][0]

PasswordResetTokenURL: [
	if #Meta.Environment.Type == "development" && #Meta.Environment.Cloud == "local" { "http://localhost:3000/he/password-reset" },
	if #Meta.Environment.Name == "staging" { "https://dev.aibookingtravel.com/he/password-reset" },
	if #Meta.Environment.Type == "production" { "https://aibookingtravel.com/he/password-reset" },
	"https://aibookingtravel.com/he/password-reset",
][0]

PriceOfferURL: [
	if #Meta.Environment.Type == "development" && #Meta.Environment.Cloud == "local" { "http://localhost:3000/he/price-offers" },
	if #Meta.Environment.Name == "staging" { "https://dev.aibookingtravel.com/he/price-offers" },
	if #Meta.Environment.Type == "production" { "https://aibookingtravel.com/he/price-offers" },
	"https://aibookingtravel.com/he/price-offers",
][0]