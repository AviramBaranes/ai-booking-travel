package accounts

import user_handlers "encore.app/services/accounts/handlers/user_handlers"

// Type aliases re-exported for backward compatibility with external callers.

type GetUserEmailParams = user_handlers.GetUserEmailParams
type GetAgentsByOfficeIDRequest = user_handlers.GetAgentsByOfficeIDParams
type GetAgentsByOrganizationIDRequest = user_handlers.GetAgentsByOrganizationIDParams
type ListAdminsEmailsResponse = user_handlers.ListAdminsEmailsResponse
type GetAgentsResponse = user_handlers.GetAgentsResponse
