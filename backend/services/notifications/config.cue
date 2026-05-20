package notifications

EmailFrom: "accounting@aibookingtravel.com"
EmailHost: "smtp.gmail.com"
EmailPort: 587

SMSUsername: "sogomatic"
SMSSenderName: "Sogomatic"

GotenbergURL: [
	if #Meta.Environment.Type == "development" && #Meta.Environment.Cloud == "local" { "http://localhost:8080" },
	if #Meta.Environment.Name == "staging" { "http://localhost:8080" },
	if #Meta.Environment.Type == "production" { "http://localhost:8080" },
	"http://localhost:8080",
][0]