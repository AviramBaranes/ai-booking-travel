package notifications

EmailFrom: "accounting@aibookingtravel.com"
EmailHost: "smtp.gmail.com"
EmailPort: 587

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