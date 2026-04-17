import { accounts } from "../client";

declare module "next-auth" {
  interface Session {
    user: accounts.LoginResponse & {
      customExp: number;
      error?: string;
      isAdminAsAgent?: boolean;
    };
  }
}
