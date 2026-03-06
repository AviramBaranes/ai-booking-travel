package auth

import "context"

type TempAgentTagUsageResponse struct {
}

// encore:api auth path=/temp-agent-tag-usage method=GET tag:agent
func (s *Service) TempAgentTagUsage(ctx context.Context) (*TempAgentTagUsageResponse, error) {
	return &TempAgentTagUsageResponse{}, nil
}

// encore:api auth path=/temp-customer-tag-usage method=GET tag:customer
func (s *Service) TempCustomerTagUsage(ctx context.Context) (*TempAgentTagUsageResponse, error) {
	return &TempAgentTagUsageResponse{}, nil
}
