import { accounts } from "../client";
import { withErrorHandler } from "./_api";

// Auth
export function login(data: accounts.LoginParams) {
  return withErrorHandler((client) => client.accounts.Login(data));
}

export function refreshTokens(data: accounts.RefreshTokensParams) {
  return withErrorHandler((client) => client.accounts.RefreshTokens(data));
}

// Admins
export function listAdmins() {
  return withErrorHandler((client) => client.accounts.ListAdmins());
}

export function createAdmin(data: accounts.CreateAdminRequest) {
  return withErrorHandler((client) => client.accounts.CreateAdmin(data));
}

// Agents
export function listAgents(data: accounts.ListAgentsRequest) {
  return withErrorHandler((client) => client.accounts.ListAgents(data));
}

export function createAgent(data: accounts.CreateAgentRequest) {
  return withErrorHandler((client) => client.accounts.CreateAgent(data));
}

// Contacts
export function listContacts(data: accounts.ListContactsRequest) {
  return withErrorHandler((client) => client.accounts.ListContacts(data));
}

export function createContact(data: accounts.CreateContactRequest) {
  return withErrorHandler((client) => client.accounts.CreateContact(data));
}

export function updateContact(id: number, data: accounts.UpdateContactRequest) {
  return withErrorHandler((client) => client.accounts.UpdateContact(id, data));
}

export function deleteContact(id: number) {
  return withErrorHandler((client) => client.accounts.DeleteContact(id));
}

// Organizations
export function listOrganizations(data: accounts.ListOrganizationsRequest) {
  return withErrorHandler((client) => client.accounts.ListOrganizations(data));
}

export function createOrganization(data: accounts.CreateOrganizationRequest) {
  return withErrorHandler((client) => client.accounts.CreateOrganization(data));
}

export function updateOrganization(
  id: number,
  data: accounts.UpdateOrganizationRequest,
) {
  return withErrorHandler((client) =>
    client.accounts.UpdateOrganization(id, data),
  );
}

// Offices
export function listOffices(data: accounts.ListOfficesRequest) {
  return withErrorHandler((client) => client.accounts.ListOffices(data));
}

export function createOffice(data: accounts.CreateOfficeRequest) {
  return withErrorHandler((client) => client.accounts.CreateOffice(data));
}

export function updateOffice(id: number, data: accounts.UpdateOfficeRequest) {
  return withErrorHandler((client) => client.accounts.UpdateOffice(id, data));
}

// Users
export function updateUser(id: number, data: accounts.UpdateUserRequest) {
  return withErrorHandler((client) => client.accounts.UpdateUser(id, data));
}
