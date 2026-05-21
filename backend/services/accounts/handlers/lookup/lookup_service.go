package lookup

import "encore.app/services/accounts/db"

type Service struct {
	query db.Querier
}

func NewService(query db.Querier) *Service {
	return &Service{query: query}
}

type AccountName struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetAccountsLookupParams struct {
	OrganizationIDs []int64 `json:"organizationIds"`
	OfficeIDs       []int64 `json:"officeIds"`
	UserIDs         []int64 `json:"userIds"`
}

type GetAccountsLookupResponse struct {
	Organizations []AccountName `json:"organizations"`
	Offices       []AccountName `json:"offices"`
	Users         []AccountName `json:"users"`
}
