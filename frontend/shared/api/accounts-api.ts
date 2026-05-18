import { accounts, auth, contact, office, organization, user } from "../client";
import { withErrorHandler } from "./_api";

// Auth
export function login(data: auth.LoginParams) {
  return withErrorHandler((client) => client.accounts.Login(data));
}

export function refreshTokens(data: auth.RefreshTokensParams) {
  return withErrorHandler((client) => client.accounts.RefreshTokens(data));
}

// Admins
export function listAdmins() {
  return withErrorHandler((client) => client.accounts.ListAdmins());
}

export function createAdmin(data: user.CreateAdminParams) {
  return withErrorHandler((client) => client.accounts.CreateAdmin(data));
}

// Accountants
export function listAccountants() {
  return withErrorHandler((client) => client.accounts.ListAccountants());
}

export function createAccountant(data: user.CreateAccountantParams) {
  return withErrorHandler((client) => client.accounts.CreateAccountant(data));
}

// Agents
export function listAgents(data: user.ListAgentsParams) {
  return withErrorHandler((client) => client.accounts.ListAgents(data));
}

export function createAgent(data: user.CreateAgentParams) {
  return withErrorHandler((client) => client.accounts.CreateAgent(data));
}

// Contacts
export function listContacts(data: contact.ListContactsParams) {
  return withErrorHandler((client) => client.accounts.ListContacts(data));
}

export function createContact(data: contact.CreateContactParams) {
  return withErrorHandler((client) => client.accounts.CreateContact(data));
}

export function updateContact(id: number, data: contact.UpdateContactParams) {
  return withErrorHandler((client) => client.accounts.UpdateContact(id, data));
}

export function deleteContact(id: number) {
  return withErrorHandler((client) => client.accounts.DeleteContact(id));
}

// Organizations
export function listOrganizations(data: organization.ListOrganizationsParams) {
  return withErrorHandler((client) => client.accounts.ListOrganizations(data));
}

export function createOrganization(data: organization.CreateOrganizationParams) {
  return withErrorHandler((client) => client.accounts.CreateOrganization(data));
}

export function updateOrganization(
  id: number,
  data: organization.UpdateOrganizationParams,
) {
  return withErrorHandler((client) =>
    client.accounts.UpdateOrganization(id, data),
  );
}

export function listOrganicOrganizations() {
  return withErrorHandler((client) => client.accounts.ListOrganicOrganizations());
}

// Offices
export function listOffices(data: office.ListOfficesParams) {
  return withErrorHandler((client) => client.accounts.ListOffices(data));
}

export function createOffice(data: office.CreateOfficeParams) {
  return withErrorHandler((client) => client.accounts.CreateOffice(data));
}

export function updateOffice(id: number, data: office.UpdateOfficeParams) {
  return withErrorHandler((client) => client.accounts.UpdateOffice(id, data));
}

export function listInorganicOffices() {
  return withErrorHandler((client) => client.accounts.ListInorganicOffices());
}

// Users
export function updateUser(id: number, data: user.UpdateUserParams) {
  return withErrorHandler((client) => client.accounts.UpdateUser(id, data));
}

export function sendOTP(data: auth.SendCustomerLoginOTPParams) {
  // Bypass withErrorHandler to avoid the 401→redirect behaviour on a public endpoint.
  return withErrorHandler((client) =>
    client.accounts.SendCustomerLoginOTP(data),
    { skipAuthRedirect: true },
  );
}

export function loginWithOTP(data: auth.ValidateCustomerLoginOTPParams) {
  return withErrorHandler((client) =>
    client.accounts.ValidateCustomerLoginOTP(data),
  );
}
 