package booking

import (
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	locations_mocks "encore.app/services/booking/mocks"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/mock/gomock"
)

func testQuerier() *db.Queries {
	pool := sqldb.Driver[*pgxpool.Pool](bookingsDB)
	return db.New(pool)
}

func mockService(t *testing.T) (*locations_mocks.MockQuerier, *Service) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := locations_mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

func invalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

func strPtr(s string) *string { return &s }
