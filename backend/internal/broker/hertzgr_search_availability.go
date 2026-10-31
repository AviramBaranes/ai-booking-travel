package broker

// SearchAvailability returns no vehicles until quoting against Autohellas is wired up. The fan-out
// reaches every broker with enabled locations, so an empty response keeps shared locations serving
// the other brokers instead of failing the whole search.
func (h *HertzGR) SearchAvailability(p SearchAvailabilityParams) (*AvailabilityResponse, error) {
	return &AvailabilityResponse{}, nil
}
