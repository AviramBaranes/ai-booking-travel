package customer

import (
	"context"
	"errors"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type ListCustomersParams struct {
	Search string `query:"search"`
	Page   int32  `query:"page" validate:"required,gte=1"`
}

func (p ListCustomersParams) Validate() error {
	return validation.ValidateStruct(p)
}

type CustomerResponse struct {
	ID          int64      `json:"id"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Email       string     `json:"email"`
	PhoneNumber *string    `json:"phoneNumber"`
	LastLogin   *time.Time `json:"lastLogin"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type ListCustomersResponse struct {
	Customers []CustomerResponse `json:"customers"`
	Total     int64              `json:"total"`
}

const PageSize int32 = 15

func (s *CustomerService) ListCustomers(ctx context.Context, p ListCustomersParams) (*ListCustomersResponse, error) {
	offset := (p.Page - 1) * PageSize

	customers, err := s.query.ListCustomers(ctx, db.ListCustomersParams{
		Search:     &p.Search,
		PageOffset: offset,
		PageSize:   PageSize,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return &ListCustomersResponse{}, nil
		}

		rlog.Error("failed to list customers", "error", err)
		return nil, api_errors.ErrInternalError
	}

	total, err := s.query.CountCustomers(ctx, &p.Search)
	if err != nil {
		rlog.Error("failed to count customers", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &ListCustomersResponse{
		Customers: toCustomerResponse(customers),
		Total:     total,
	}, nil
}

func toCustomerResponse(rows []db.ListCustomersRow) []CustomerResponse {
	customers := make([]CustomerResponse, len(rows))
	for i, r := range rows {
		lastLogin := dbadapters.TimeFromDB(r.LastLogin)
		customers[i] = CustomerResponse{
			ID:          r.ID,
			FirstName:   r.FirstName,
			LastName:    r.LastName,
			Email:       r.Email,
			PhoneNumber: r.PhoneNumber,
			LastLogin:   &lastLogin,
			CreatedAt:   dbadapters.TimeFromDB(r.CreatedAt),
			UpdatedAt:   dbadapters.TimeFromDB(r.UpdatedAt),
		}
	}
	return customers
}
