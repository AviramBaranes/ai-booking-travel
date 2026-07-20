package customer

import "encore.app/services/accounts/db"

type CustomerService struct {
	query db.Querier
}

func NewCustomerService(query db.Querier) *CustomerService {
	return &CustomerService{
		query: query,
	}
}
