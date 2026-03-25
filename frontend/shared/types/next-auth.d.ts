import { auth } from "../client";

declare module "next-auth" {
  interface Session {
    user: auth.LoginResponse & { customExp: number };
  }
}
