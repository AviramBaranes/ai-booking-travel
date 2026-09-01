package broker

// SearchAvailability returns no vehicles until quoting against Wheelsys is wired up. The fan-out
// reaches every broker with enabled locations, so an empty response keeps shared locations serving
// the other brokers instead of failing the whole search.
func (a *Avance) SearchAvailability(p SearchAvailabilityParams) (*AvailabilityResponse, error) {
	return &AvailabilityResponse{}, nil
}
