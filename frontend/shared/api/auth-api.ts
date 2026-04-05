import { accounts } from "../client";
import { withErrorHandler } from "./_api";

export async function login(data: accounts.LoginParams) {
  return withErrorHandler((client) => client.accounts.Login(data));
}
