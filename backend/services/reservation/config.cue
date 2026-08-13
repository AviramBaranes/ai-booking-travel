package reservation

VAT: 18.0

GotenbergURL: [
	if #Meta.Environment.Type == "development" && #Meta.Environment.Cloud == "local" { "http://localhost:8080" },
	if #Meta.Environment.Name == "staging" { "https://gotenberg-217822056800.me-west1.run.app" },
	if #Meta.Environment.Type == "production" { "https://gotenberg-217822056800.me-west1.run.app" },
	"https://gotenberg-217822056800.me-west1.run.app",
][0]

Icount:{
	CID: "aibookingtravel",
	User: "accounting",
	// iCount supplier ids, per broker, used when recording what we paid the supplier.
	SupplierIDs: {
		Flex: 5,
	}
	// The iCount expense type supplier payments are recorded under ("תשלום רכבים").
	ExpenseTypeID: 5,
}