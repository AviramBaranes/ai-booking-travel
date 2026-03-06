package supplier

type SupplierName string

const (
	SupplierFlex  SupplierName = "flex"
	SupplierHertz SupplierName = "hertz"
)

type Supplier interface {
	Supplier() SupplierName
	GetLocationsPage(cursor string) (LocationPage, error)
}

type LocationPage struct {
	Locations []Location
	NextPage  string
}

type Location struct {
	ID          string
	Name        string
	Country     string
	City        string
	CountryCode string
	Iata        string
}
