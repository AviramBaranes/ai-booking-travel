import { auth } from "../client";
import { withErrorHandler } from "./_api";

export async function login(data: auth.LoginParams) {
  return withErrorHandler((client) => client.auth.Login(data));
}
