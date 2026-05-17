package booking

import (
	"context"

	availability "encore.app/services/booking/handlers/availability"
	"encore.dev/config"
)

// AvCfg is the loaded AvailableVehiclesConfig for this service.
var AvCfg = config.Load[*availability.AvailableVehiclesConfig]()

// SearchAvailability handles the http request for searching availability of vehicles.
// encore:api public method=GET path=/booking/availability
func (s *Service) SearchAvailability(ctx context.Context, p availability.SearchAvailabilityParams) (*availability.SearchAvailabilityResponse, error) {
	as := availability.NewAvailabilityService(s.query, s.t, AvCfg)
	return as.SearchAvailability(ctx, p)
}
