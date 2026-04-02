package accounts

import "context"

type TempAgentTagUsageResponse struct {
}

// encore:api auth path=/temp-customer-tag-usage method=GET tag:customer
func (s *Service) TempCustomerTagUsage(ctx context.Context) (*TempAgentTagUsageResponse, error) {
	return &TempAgentTagUsageResponse{}, nil
}
